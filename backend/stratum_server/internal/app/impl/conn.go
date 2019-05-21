package impl

import (
	"github.com/davyxu/cellnet"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/context"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/app/server"
	"net"
	"time"
	"github.com/prep/average"
	log "github.com/sirupsen/logrus"
)

func Accepted(s *server.Server, ev cellnet.Event) {
	contextSet := ev.Session().Peer().(cellnet.ContextSet)
	stratumContext := context.NewStratumContext()
	sid := ev.Session().ID()
	remoteIP := ev.Session().Raw().(*net.TCPConn).RemoteAddr()
	localIP := ev.Session().Raw().(*net.TCPConn).LocalAddr()
	difficultyCheckInterval := s.GetConfig().StratumConfig.DifficultyCheckLoopInterval

	windowDuration := difficultyCheckInterval * time.Second
	slideWindow := average.MustNew(5*windowDuration, windowDuration)

	stratumContext.LocalIP = localIP.String()
	stratumContext.RemoteIP = remoteIP.String()
	stratumContext.AcceptedTs = time.Now().Unix()
	stratumContext.SlideWindow = slideWindow
	contextSet.SetContext(sid, stratumContext)
	log.WithFields(log.Fields{
		"sid":       sid,
		"local_ip":  localIP,
		"remote_ip": remoteIP,
	}).Debugln("session accepted.")
}

func Closed(s *server.Server, ev cellnet.Event) {
	contextSet := ev.Session().Peer().(cellnet.ContextSet)
	sid := ev.Session().ID()
	stratumContext, ok := contextSet.GetContext(sid)
	if ok {
		stratumContext.(*context.StratumContext).ClosedTs = time.Now().Unix()
		stratumContext.(*context.StratumContext).CancelFunc()
		stratumContext.(*context.StratumContext).SlideWindow.Stop()
		contextSet.SetContext(sid, nil)
	}

	log.WithFields(log.Fields{
		"sid": sid,
	}).Debugln("session closed.")
}

func Unknown(s *server.Server, ev cellnet.Event) {
	log.WithFields(log.Fields{
		"sid":     ev.Session().ID(),
		"message": ev.Message(),
	}).Debugln("session method unknown.")
}
