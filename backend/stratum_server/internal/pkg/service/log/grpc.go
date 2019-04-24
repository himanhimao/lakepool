package log

import (
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/conf"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_log"
	"fmt"
	"google.golang.org/grpc"
	"time"
	"context"
)

const (
	Timeout = time.Second * 30
)

type GRPCService struct {
	config     *conf.GRPCConfig
	client     pb.LogClient
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
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	s.client = pb.NewLogClient(conn)
	return nil
}

func (s *GRPCService) AddBlockLog(coinType string, t time.Time, blockLog *service.BlockLog) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	addBlockLogRequest := &pb.AddBlockLogRequest{
		Ts:       t.UnixNano(),
		CoinType: coinType,
		Log:      loadPbBlockLog(blockLog),
	}

	addShareLogResponse, err := s.client.AddBlockLog(ctx, addBlockLogRequest)
	if err != nil {
		return false, err
	}

	return addShareLogResponse.Result, nil
}

func loadPbBlockLog(log *service.BlockLog) *pb.BlockLog {
	pbBlockLog := new(pb.BlockLog)
	pbBlockLog.HostName = log.GetHostName()
	pbBlockLog.UserAgent = log.GetUserAgent()
	pbBlockLog.WorkerName = log.GetWorkerName()
	pbBlockLog.Height = log.GetHeight()
	pbBlockLog.Pid = log.GetPid()
	pbBlockLog.ExtName = log.GetExtName()
	pbBlockLog.UserName = log.GetUserName()
	pbBlockLog.ClientIp = log.GetClientIP()
	pbBlockLog.ServerIp = log.GetServerIP()
	pbBlockLog.Hash = log.GetHash()
	return pbBlockLog
}
