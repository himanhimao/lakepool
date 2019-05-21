package impl

import (
	"encoding/hex"
	"github.com/davyxu/cellnet"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/context"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/app/server"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/util"
	"strings"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/peer/tcp"
	ctx "context"
)

func Subscribe(s *server.Server, ev cellnet.Event) {
	subScribeREQ := proto.NewSubscribeREQ()
	subScribeREQ.Load(ev.Message().(*proto.JSONRpcREQ))

	contextSet := ev.Session().Peer().(cellnet.ContextSet)
	sid := ev.Session().ID()
	msgId := subScribeREQ.Id
	config := s.GetConfig().StratumConfig
	extraNonce1 := hex.EncodeToString([]byte(util.RandString(config.StratumSubscribeConfig.ExtraNonce1Length)))
	extraNonce2LengthValue := config.StratumSubscribeConfig.ExtraNonce2Length
	notifyPlaceHolder := config.StratumSubscribeConfig.DifficultyPlaceholder
	difficultyPlaceHolder := config.StratumSubscribeConfig.DifficultyPlaceholder
	userAgent := subScribeREQ.UserAgent

	stratumContext, ok := contextSet.GetContext(sid)
	if ok {
		if len(userAgent) > 0 {
			stratumContext.(*context.StratumContext).UserAgent = userAgent
		}
		stratumContext.(*context.StratumContext).SubscribeTs = time.Now().Unix()
		stratumContext.(*context.StratumContext).SessionID = extraNonce1
	} else {
		log.WithFields(log.Fields{
			"sid": sid,
		}).Errorln("not found context")
		errRESP := proto.NewErrOtherUnknownRESP(msgId)
		ev.Session().Send(errRESP)
		ev.Session().Close()
		return
	}

	log.WithFields(log.Fields{
		"sid":                       sid,
		"msg_id":                    msgId,
		"user_agent":                userAgent,
		"extra_nonce_1":             extraNonce1,
		"extra_nonce2_length_value": extraNonce2LengthValue,
		"notify":                    notifyPlaceHolder,
		"difficulty":                difficultyPlaceHolder,
	}).Infoln("subscribed.")

	resp := proto.NewSubscribeRESP(msgId, notifyPlaceHolder, difficultyPlaceHolder, extraNonce1, extraNonce2LengthValue)
	ev.Session().Send(resp)
}

func Authorize(s *server.Server, ev cellnet.Event) {
	authorizeREQ := proto.NewAuthorizeREQ()
	authorizeREQ.Load(ev.Message().(*proto.JSONRpcREQ))

	contextSet := ev.Session().Peer().(cellnet.ContextSet)
	sid := ev.Session().ID()
	msgId := authorizeREQ.Id
	workerName := strings.TrimSpace(authorizeREQ.WorkerName)
	config := s.GetConfig().StratumConfig
	stratumContext, ok := contextSet.GetContext(sid)

	errUnknownRESP := proto.NewErrOtherUnknownRESP(msgId)
	errNotSubscribedRESP := proto.NewErrNotSubscribedRESP(msgId)

	if !ok {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Errorln("not found context")
		ev.Session().Send(errUnknownRESP)
		ev.Session().Close()
		return
	}

	if stratumContext.(*context.StratumContext).SubscribeTs == 0 {
		log.WithFields(log.Fields{
			"sid": sid,
		}).Errorln("not Subscribed")
		ev.Session().Send(errNotSubscribedRESP)
		ev.Session().Close()
		return
	}

	if len(workerName) == 0 {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Errorln("full name is empty")
		ev.Session().Send(errUnknownRESP)
		ev.Session().Close()
		return
	}

	var username, extName string
	nameSeparator := config.NameSeparator
	nameSlice := strings.SplitN(workerName, nameSeparator, 2)
	if len(nameSlice) == 2 {
		username = strings.TrimSpace(nameSlice[0])
		extName = strings.TrimSpace(nameSlice[1])
	} else {
		username = workerName
	}

	password := authorizeREQ.Password
	stratumContext.(*context.StratumContext).WorkerName = workerName
	stratumContext.(*context.StratumContext).UserName = username
	stratumContext.(*context.StratumContext).ExtName = extName
	stratumContext.(*context.StratumContext).Password = password
	stratumContext.(*context.StratumContext).AcceptedTs = time.Now().Unix()

	log.WithFields(log.Fields{
		"sid":         sid,
		"msg_id":      msgId,
		"worker_name": workerName,
		"user_name":   username,
		"ext_name":    extName,
		"password":    password,
	}).Infoln("authorized")

	var result proto.AuthorizeResult

	result = s.GetServiceMgr().GetUserService().(service.UserService).Login(workerName, password)
	resp := proto.NewAuthorizeRESP(msgId, result)
	ev.Session().Send(resp)

	log.WithFields(log.Fields{
		"sid":         sid,
		"worker_name": workerName,
	}).Infoln("user authorize result:", result)

	//验证不通过 返回
	if result != proto.AuthorizePass {
		return
	}

	//验证通过 发送难度数据 还有任务数据
	var defaultDifficulty uint64
	userAgent := stratumContext.(*context.StratumContext).UserAgent

	if len(userAgent) == 0 {
		defaultDifficulty = config.DefaultDifficulty
	} else {
		if difficulty, err := config.GetNotifyDefaultDifficulty(userAgent); err != nil {
			defaultDifficulty = config.DefaultDifficulty
		} else {
			defaultDifficulty = difficulty
		}
	}

	//发送难度
	stratumContext.(*context.StratumContext).LatestDifficulty = defaultDifficulty
	difficultyResp := proto.NewJSONRpcSetDifficultyRESP(defaultDifficulty)
	ev.Session().Send(difficultyResp)

	//发送任务
	latestJobHeight := s.GetJobRepo().GetLatestHeight()
	latestJobIndex := 0
	latestNotifyTs := time.Now().Unix()
	startNotifyTs := latestNotifyTs
	jobId := service.GenerateJobId(latestJobHeight, latestJobIndex, defaultDifficulty)
	stratumJob := s.GetJobRepo().GetJob(latestJobHeight, latestJobIndex)

	if stratumJob == nil {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Errorln("not found job")
		ev.Session().Send(errUnknownRESP)
		ev.Session().Close()
		return
	}

	stratumJob = stratumJob.Fill(jobId, latestNotifyTs, true)
	stratumContext.(*context.StratumContext).LatestNotifyJobHeight = latestJobHeight
	stratumContext.(*context.StratumContext).StartNotifyTs = startNotifyTs
	stratumContext.(*context.StratumContext).LatestNotifyTs = latestNotifyTs
	stratumContext.(*context.StratumContext).LatestNotifyJobIndex = latestJobIndex
	stratumContext.(*context.StratumContext).NotifyCount ++

	notifyRESP := proto.NewJSONRpcNotifyRESP(stratumJob.ToJSONInterface())
	ev.Session().Send(notifyRESP)

	//任务检查轮询&难度检查轮询&下发任务轮询
	go func(ctx ctx.Context, ses *tcp.StratumSession) {
		log.Debugln(sid, workerName, "goroutine start")
		notifyLoopTicker := time.NewTicker(time.Second * s.GetConfig().StratumConfig.NotifyLoopInterval)
		difficultyCheckTicker := time.NewTicker(time.Second * s.GetConfig().StratumConfig.DifficultyCheckLoopInterval)
		jobCheckInterval := time.NewTicker(time.Millisecond * s.GetConfig().StratumConfig.JobHashCheckLoopInterval)
		for {
			select {
			case <-ctx.Done():
				log.WithFields(log.Fields{
					"sid":         sid,
					"worker_name": workerName,
				}).Debugln("goroutine close")
				return
			case <-notifyLoopTicker.C:
				ses.ProcEvent(&cellnet.RecvMsgEvent{Ses: ses, Msg: &proto.JobNotify{}})
			case <-difficultyCheckTicker.C:
				ses.ProcEvent(&cellnet.RecvMsgEvent{Ses: ses, Msg: &proto.JobDifficultyCheck{}})
			case <-jobCheckInterval.C:
				ses.ProcEvent(&cellnet.RecvMsgEvent{Ses: ses, Msg: &proto.JobHeightCheck{}})
			}
		}
	}(stratumContext.(*context.StratumContext).Context, ev.Session().(*tcp.StratumSession))
}

func Submit(s *server.Server, ev cellnet.Event) {
	submitREQ := proto.NewSubmitREQ()
	submitREQ.Load(ev.Message().(*proto.JSONRpcREQ))

	contextSet := ev.Session().Peer().(cellnet.ContextSet)
	sid := ev.Session().ID()
	msgId := submitREQ.Id
	workerName := submitREQ.WorkerName
	jobId := submitREQ.JobId
	nTime := submitREQ.NTime
	nonce := submitREQ.Nonce
	extraNonce2 := submitREQ.ExtraNonce2
	log.WithFields(log.Fields{
		"sid":          sid,
		"msg_id":       msgId,
		"worker_name":  workerName,
		"job_id":       jobId,
		"nonce":        nonce,
		"extra_nonce2": extraNonce2,
		"n_time":       nTime,
	}).Infoln("submitted")

	stratumContext, ok := contextSet.GetContext(sid)
	errUnknownRESP := proto.NewErrOtherUnknownRESP(msgId)
	errUnauthorizedWorkerResp := proto.NewErrUnauthorizedWorkerRESP(msgId)
	stratumContext.(*context.StratumContext).SubmitCount++

	if stratumContext.(*context.StratumContext).AuthorizeTs == 0 || workerName != stratumContext.(*context.StratumContext).WorkerName {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Errorln("unauthorized worker.")
		ev.Session().Send(errUnauthorizedWorkerResp)
		return
	}

	if !ok {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Errorln("not found context")
		ev.Session().Send(errUnknownRESP)
		ev.Session().Close()
		return
	}

	height, index, difficulty, err := service.ExtractJobId(jobId)
	if err != nil {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
			"error":       err,
		}).Errorln("invalid jobId")
		ev.Session().Send(errUnknownRESP)
		ev.Session().Close()
		return
	}
	log.WithFields(log.Fields{
		"sid":        sid,
		"height":     height,
		"index":      index,
		"difficulty": difficulty,
	}).Debugln("submitted info")

	errJobNotFoundResp := proto.NewErrJobNotFoundRESP(msgId)
	job := s.GetJobRepo().GetJob(height, index)
	nTimeTs, err := util.ConvertTs(nTime)
	if err != nil {
		log.WithFields(log.Fields{
			"sid":    sid,
			"n_time": nTime,
			"error":  err,
		}).Errorln("invalid nTime, conversion error")
		ev.Session().Send(errUnknownRESP)
		ev.Session().Close()
		return
	}

	if s.GetJobRepo().GetLatestHeight() != height || job == nil || uint32(job.GetMeta().GetMinTimeTs()) > nTimeTs {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Errorln("stratum servce job not found")
		ev.Session().Send(errJobNotFoundResp)
		ev.Session().Close()
		return
	}

	slideWindow := stratumContext.(*context.StratumContext).SlideWindow
	if slideWindow == nil {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Errorln("not found slide window")
		ev.Session().Send(errUnknownRESP)
		ev.Session().Close()
		return
	}
	slideWindow.Add(1)

	extraNonce1 := stratumContext.(*context.StratumContext).SessionID
	share := job.ToShare()
	share.SetExtraNonce1(extraNonce1).SetExtraNonce2(extraNonce2).SetNTime(nTime).SetNonce(nonce)

	shareResult, err := s.GetServiceMgr().GetSphereService().SubmitShare(share, difficulty)
	submitRefuseRESP := proto.NewSubmitRESP(msgId, proto.SubmitRefuse)
	if err != nil {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
			"error":       err,
		}).Errorln("submit share error")
		stratumContext.(*context.StratumContext).ErrorCount++
		ev.Session().Send(submitRefuseRESP)
		return
	}

	var isRightShare bool
	if shareResult.State != service.StateSuccess && shareResult.State != service.StateSuccessSubmitBlock {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
			"status":      shareResult.State,
		}).Infoln("submit share failed")
		stratumContext.(*context.StratumContext).RejectCount++
		ev.Session().Send(submitRefuseRESP)
	} else {
		isRightShare = true
		log.WithFields(log.Fields{
			"sid":           sid,
			"worker_name":   workerName,
			"compute_power": shareResult.ComputePower,
		}).Infoln("submit share success")
		submitPassRESP := proto.NewSubmitRESP(msgId, proto.SubmitPass)
		stratumContext.(*context.StratumContext).AcceptCount++
		ev.Session().Send(submitPassRESP)
	}

	//add share log
	hostName, _ := s.GetSysInfo().GetHostName()
	pid := int32(s.GetSysInfo().GetPid())
	clientIp := stratumContext.(*context.StratumContext).RemoteIP
	serverIp := stratumContext.(*context.StratumContext).LocalIP
	userAgent := stratumContext.(*context.StratumContext).UserAgent
	username := stratumContext.(*context.StratumContext).UserName
	extName := stratumContext.(*context.StratumContext).ExtName
	coinType := s.GetConfig().StratumConfig.CoinType
	t := time.Now()

	shareLog := service.NewShareLog().SetHeight(height).SetClientIP(clientIp).
		SetServerIP(serverIp).SetHostName(hostName).SetPid(pid).SetWorkerName(workerName).SetUserName(username).
		SetExtName(extName).SetUserAgent(userAgent).SetIsRight(isRightShare).SetComputePower(float64(shareResult.ComputePower))
	_, err = s.GetServiceMgr().GetStatsService().AddShareLog(coinType, t, shareLog)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorln("add share log fail")
	} else {
		log.Infoln("add share log success")
	}

	//add block log
	if shareResult.State == service.StateSuccessSubmitBlock {
		blockLog := service.NewBlockLog().SetHeight(height).SetClientIP(clientIp).
			SetServerIP(serverIp).SetHostName(hostName).SetPid(pid).SetWorkerName(workerName).SetUserName(username).
			SetExtName(extName).SetUserAgent(userAgent).SetHash(string(shareResult.BlockHash))
		_, err = s.GetServiceMgr().GetLogService().AddBlockLog(coinType, t, blockLog)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Errorln("add block log fail")
		} else {
			log.Infoln("add block log success")
		}
	}
	return
}