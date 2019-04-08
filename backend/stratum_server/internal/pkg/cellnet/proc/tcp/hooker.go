package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"
	log "github.com/sirupsen/logrus"

)

// 带有RPC和relay功能
type MsgHooker struct {
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {
	switch inputEvent.Message().(type) {
	case *proto.JobDifficultyCheck:
		return inputEvent
	case *proto.JobHeightCheck:
		return inputEvent
	case *proto.JobNotify:
		return inputEvent
	}
	log.WithFields(log.Fields{
		"input_event_session":inputEvent.Session(),
		"input_event_message": inputEvent.Message(),
	}).Debugln("tcp.Stratum")
	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {
	log.WithFields(log.Fields{
		"output_event_session":inputEvent.Session(),
		"output_event_message": inputEvent.Message(),
	}).Debugln("tcp.Stratum")
	return inputEvent
}
