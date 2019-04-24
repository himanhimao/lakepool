package stats

import (
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/conf"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_stats"
	"google.golang.org/grpc"
	"time"
	"context"
)

const (
	Timeout = time.Second * 30
)

type GRPCService struct {
	config     *conf.GRPCConfig
	client     pb.StatsClient
	registerId string
}

func NewGRPCService() *GRPCService {
	return &GRPCService{}
}

func (s *GRPCService) SetConfig(config *conf.GRPCConfig) *GRPCService {
	s.config = config
	return s
}

func (s *GRPCService) Init() error {
	address := s.config.FormatHostPort()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	s.client = pb.NewStatsClient(conn)
	return nil
}

func (s *GRPCService) AddShareLog(coinType string, t time.Time, shareLog *service.ShareLog) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	addShareLogRequest := &pb.AddShareLogRequest{
		Ts:       t.UnixNano(),
		CoinType: coinType,
		Log:      loadPbShareLog(shareLog),
	}

	addShareLogResponse, err := s.client.AddShareLog(ctx, addShareLogRequest)
	if err != nil {
		return false, err
	}

	return addShareLogResponse.Result, nil
}

func loadPbShareLog(log *service.ShareLog) *pb.ShareLog {
	pbShareLog := new(pb.ShareLog)
	pbShareLog.HostName = log.GetHostName()
	pbShareLog.UserAgent = log.GetUserAgent()
	pbShareLog.WorkerName = log.GetWorkerName()
	pbShareLog.Height = log.GetHeight()
	pbShareLog.Pid = log.GetPid()
	pbShareLog.ExtName = log.GetExtName()
	pbShareLog.UserName = log.GetUserName()
	pbShareLog.ClientIp = log.GetClientIP()
	pbShareLog.ServerIp = log.GetServerIP()
	pbShareLog.IsRight = log.GetIsRight()
	pbShareLog.ComputePower = log.GetComputePower()
	return pbShareLog
}
