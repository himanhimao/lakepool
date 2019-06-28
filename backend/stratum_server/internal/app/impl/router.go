package impl

import (
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/app/server"
)

func RegisterRoute(router *server.Router) {
	router.RegisterRouter(proto.MiningRecvMethodSubscribe, Subscribe)
	router.RegisterRouter(proto.MiningRecvMethodExtranonceSubscribe, ExtraSubscribe)
	router.RegisterRouter(proto.MiningRecvMethodAuthorize, Authorize)
	router.RegisterRouter(proto.MiningRecvMethodSubmit, Submit)
	router.RegisterRouter(proto.ConnClosed, Closed)
	router.RegisterRouter(proto.ConnAcceptd, Accepted)
	router.RegisterRouter(proto.ConnUnknown, Unknown)
	router.RegisterRouter(proto.JobMethodNotify, JobNotify)
	router.RegisterRouter(proto.JobMethodDifficultyCheck, JobDifficultyCheck)
	router.RegisterRouter(proto.JobMethodHeightCheck, JobHeightCheck)
}
