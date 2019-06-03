package server

import (
	"context"
	"encoding/json"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/json"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/peer/tcp"
	_ "github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proc/tcp"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"
	_ "github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/conf"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type RouterHandler func(*Server, cellnet.Event)

type Router struct {
	routeTable map[string]RouterHandler
}

type Server struct {
	config     *conf.Config
	sMgr       *service.Manager
	sysInfo    *service.SysInfo
	jobRepo    *JobRepo
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (s *Server) GetConfig() *conf.Config {
	return s.config
}

func (s *Server) GetServiceMgr() *service.Manager {
	return s.sMgr
}

func (s *Server) GetJobRepo() *JobRepo {
	return s.jobRepo
}

func (s *Server) GetSysInfo() *service.SysInfo {
	return s.sysInfo
}

func NewRouter(size int) *Router {
	var table map[string]RouterHandler
	if size > 0 {
		table = make(map[string]RouterHandler, size)
	} else {
		table = make(map[string]RouterHandler)
	}
	return &Router{routeTable: table}
}

func (s *Router) RegisterRouter(name string, handler RouterHandler) {
	s.routeTable[name] = handler
}

func (s *Router) getHandler(name string) RouterHandler {
	return s.routeTable[name]
}

func NewServer() *Server {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &Server{jobRepo: NewJobRepo(), ctx: ctx, cancelFunc: cancel}
}

func (s *Server) Cancel() {
	s.cancelFunc()
}

func (s *Server) SetConfig(config *conf.Config) *Server {
	s.config = config
	return s
}

func (s *Server) SetServiceMgr(sMgr *service.Manager) *Server {
	s.sMgr = sMgr
	return s
}

func (s *Server) Init() error {
	var err error
	if err = s.config.Init(); err != nil {
		return err
	}

	if err = s.sMgr.Init(); err != nil {
		return err
	}

	if err = s.registerSphere(); err != nil {
		return err
	}

	if err = s.loadNewJob(); err != nil {
		return err
	}

	go s.sMgr.GetSphereService().Subscribe(s.ctx, func(job *service.StratumJob) {
		if s.jobRepo.GetLatestHeight() != job.Meta.Height {
			if _, err := s.sMgr.GetSphereService().ClearShareHistory(s.jobRepo.GetLatestHeight()); err != nil {
				log.Errorln("clear share history error.", err)
			}
		}
		index := s.jobRepo.SetJob(job.Meta.Height, job)
		log.WithFields(log.Fields{
			"height": job.Meta.Height,
			"index":  index,
			"timestamp": job.Meta.CurTimeTs,
		}).Infoln("subscribe new job")
	})
	return nil
}

func (s *Server) registerSphere() error {
	var err error
	if err = s.sMgr.GetSphereService().Init(); err != nil {
		return err
	}

	serviceStratumConfig := service.NewStratumConfig()
	config := s.config.StratumConfig
	serviceStratumConfig.PoolTag = config.PoolTag
	serviceStratumConfig.PayoutAddress = config.PayoutAddress
	serviceStratumConfig.CoinType = config.CoinType
	serviceStratumConfig.ExtraNonce2Length = config.StratumSubscribeConfig.ExtraNonce2Length
	serviceStratumConfig.ExtraNonce1Length = config.StratumSubscribeConfig.ExtraNonce1Length

	s.sysInfo = service.NewSysInfo()
	if err = s.sMgr.GetSphereService().Register(serviceStratumConfig, s.sysInfo); err != nil {
		return err
	}

	return err
}

func (s *Server) loadNewJob() error {
	var err error
	var stratumJob *service.StratumJob
	jobRepo := s.jobRepo
	sphereService := s.sMgr.GetSphereService()
	stratumJob, err = sphereService.GetLatestStratumJob()
	if err != nil {
		return err
	}
	jobRepo.SetJob(stratumJob.Meta.Height, stratumJob)
	log.Debugln("init job height", stratumJob.Meta.Height)
	go jobRepo.Clean(s.ctx)
	return err
}

func (s *Server) Run(router *Router) {
	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	queue := cellnet.NewEventQueue()

	// 创建一个tcp的侦听器，名称为server，连接地址为127.0.0.1:8801，所有连接将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	p := peer.NewGenericPeer("tcp.Stratum.Acceptor", s.config.ServerConfig.Name, s.config.ServerConfig.FormatHostPort(), queue)

	// 每一个连接收到的所有消息事件(cellnet.Event)都被派发到用户回调, 用户使用switch判断消息类型，并做出不同的处理
	proc.BindProcessorHandler(p, "tcp.stratum", func(ev cellnet.Event) {
		var handler RouterHandler
		var msgId interface{}
		switch msg := ev.Message().(type) {
		case *cellnet.SessionAccepted:
			handler = router.getHandler(proto.ConnAcceptd)
			// 有连接断开
		case *cellnet.SessionClosed:
			// 收到某个连接的ChatREQ消息
			handler = router.getHandler(proto.ConnClosed)
		case *proto.JSONRpcREQ:
			log.Debugln("method:", msg.Method)
			handler = router.getHandler(msg.Method)
			msgId = msg.Id
		case *proto.JobNotify:
			handler = router.getHandler(proto.JobMethodNotify)
		case *proto.JobDifficultyCheck:
			handler = router.getHandler(proto.JobMethodDifficultyCheck)
		case *proto.JobHeightCheck:
			handler = router.getHandler(proto.JobMethodHeightCheck)
		default:
			handler = router.getHandler(proto.ConnUnknown)
		}

		if handler != nil {
			handler(s, ev)
		} else {
			if msgId != nil {
				otherUnknownRESP := proto.NewErrOtherUnknownRESP(msgId.(*json.RawMessage))
				ev.Session().Send(otherUnknownRESP)
				ev.Session().Close()
			} else {
				ev.Session().Close()
			}
		}
	})
	// 开始侦听
	p.Start()

	// 事件队列开始循环
	queue.StartLoop()

	//信号捕捉
	sigalChannel := make(chan os.Signal, 1)
	signal.Notify(sigalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	go func(queue cellnet.EventQueue, c chan os.Signal) {
		sig := <-c
		log.Infoln(sig, "server shutting down......")
		s.Cancel()
		if _, err := s.sMgr.GetSphereService().UnRegister(); err != nil {
			log.Warnln("unregister error.", err)
		}
		queue.StopLoop()
	}(queue, sigalChannel)

	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	queue.Wait()
}
