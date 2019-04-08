package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

// 接受器
type StratumAcceptor struct {
	peer.SessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle
	peer.CoreTCPSocketOption
	peer.CoreCaptureIOPanic

	// 保存侦听器
	listener net.Listener
}

func (self *StratumAcceptor) Port() int {
	if self.listener == nil {
		return 0
	}

	return self.listener.Addr().(*net.TCPAddr).Port
}

func (self *StratumAcceptor) IsReady() bool {
	return self.IsRunning()
}

//// 异步开始侦听
func (self *StratumAcceptor) Start() cellnet.Peer {
	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	ln, err := util.DetectPort(self.Address(), func(a *util.Address, port int) (interface{}, error) {
		return net.Listen("tcp", a.HostPortString(port))
	})

	if err != nil {

		log.Errorf("#tcp.listen failed(%s) %v", self.Name(), err.Error())

		self.SetRunning(false)

		return self
	}

	self.listener = ln.(net.Listener)

	log.Infof("#tcp.listen(%s) %s", self.Name(), self.ListenAddress())

	go self.accept()

	return self
}

func (self *StratumAcceptor) ListenAddress() string {
	pos := strings.Index(self.Address(), ":")
	if pos == -1 {
		return self.Address()
	}

	host := self.Address()[:pos]

	return util.JoinAddress(host, self.Port())
}

func (self *StratumAcceptor) accept() {
	self.SetRunning(true)

	for {
		conn, err := self.listener.Accept()

		if self.IsStopping() {
			break
		}

		if err != nil {
			// 调试状态时, 才打出accept的具体错误
			log.Debugln("#tcp.accept failed(%s) %v", self.Name(), err.Error())
			break
		}

		// 处理连接进入独立线程, 防止accept无法响应
		go self.onNewSession(conn)
	}

	self.SetRunning(false)
	self.EndStopping()
}

func (self *StratumAcceptor) onNewSession(conn net.Conn) {

	self.ApplySocketOption(conn)

	ses := newSession(conn, self, nil)

	ses.Start()
	self.ProcEvent(&cellnet.RecvMsgEvent{
		Ses: ses,
		Msg: &cellnet.SessionAccepted{},
	})
}

// 停止侦听器
func (self *StratumAcceptor) Stop() {
	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()
	self.listener.Close()
	// 断开所有连接
	self.CloseAllSession()
	// 等待线程结束
	self.WaitStopFinished()
}

func (self *StratumAcceptor) TypeName() string {
	return "tcp.Stratum.Acceptor"
}

func init() {
	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &StratumAcceptor{
			SessionManager: new(peer.CoreSessionManager),
		}
		p.CoreTCPSocketOption.Init()

		return p
	})
}
