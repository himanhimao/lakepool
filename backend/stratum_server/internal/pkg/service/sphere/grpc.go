package sphere

import (
	"context"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/conf"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_sphere"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/connectivity"
	"io"
	"time"
)

var (
	RetryNum          = 1000
	ReconnectWaitTime = time.Second * 15
	Timeout           = time.Second * 30
)

type GRPCService struct {
	config     *conf.GRPCConfig
	client     pb.SphereClient
	conn       *grpc.ClientConn
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
	return s.connect()
}

func (s *GRPCService) connect() error {
	address := s.config.FormatHostPort()
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		return err
	}
	s.client = pb.NewSphereClient(conn)
	s.conn = conn
	return nil
}

func (s *GRPCService) reconnect() error {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return err
		}
	}
	return s.connect()
}

func (s *GRPCService) GetLatestStratumJob() (*service.StratumJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	getLatestStratumJobRequest := &pb.GetLatestStratumJobRequest{RegisterId: s.registerId}
	getLatestStratumJobResponse, err := s.client.GetLatestStratumJob(ctx, getLatestStratumJobRequest)

	if err != nil {
		return nil, err
	}

	pbStratumJob := getLatestStratumJobResponse.Job
	//远程获取
	return loadStratumJob(pbStratumJob), nil
}

func (s *GRPCService) Register(config *service.StratumConfig, sysInfo *service.SysInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	pbStratumConfig := new(pb.StratumConfig)
	pbStratumConfig.CoinType = config.CoinType
	pbStratumConfig.PayoutAddress = config.PayoutAddress
	pbStratumConfig.PoolTag = config.PoolTag
	pbStratumConfig.ExtraNonce1Length = int32(config.ExtraNonce1Length)
	pbStratumConfig.ExtraNonce2Length = int32(config.ExtraNonce2Length)
	pbSysInfo := new(pb.SysInfo)
	pbSysInfo.Hostname, _ = sysInfo.GetHostName()
	pbSysInfo.Pid = int32(sysInfo.GetPid())

	registerRequest := &pb.RegisterRequest{Config: pbStratumConfig, SysInfo: pbSysInfo}
	registerResponse, err := s.client.Register(ctx, registerRequest)
	if err != nil {
		return err
	}
	s.registerId = registerResponse.RegisterId
	return nil
}

func (s *GRPCService) Subscribe(ctx context.Context, handler service.SubscribeHandler) {
	var retryNum int
	var stream pb.Sphere_SubscribeClient
	var err error
	log.Infoln("Subscribing to gbt job")
Retry:
	if s.conn.GetState() != connectivity.Ready || stream == nil {
		log.Infoln("subscribe connect")
		time.Sleep(ReconnectWaitTime)
		err = s.reconnect()
		retryNum++
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Errorln("subscribe connect error")
			if retryNum < RetryNum {
				goto Retry
			} else {
				goto Out
			}
		}

		getLatestStratumJobRequest := &pb.GetLatestStratumJobRequest{RegisterId: s.registerId}
		stream, err = s.client.Subscribe(context.Background(), getLatestStratumJobRequest)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Errorln("subscribe remote method error")
			goto Retry
		}
	}

	for {
		if stream != nil {
			select {
			case <-stream.Context().Done():
				log.Errorln("sphere subscribe server interrupt. try reconnect..")
				goto Retry
			case <-ctx.Done():
				log.Errorln("subscribe routine cancel.")
				goto Out
			default:
				resp, err := stream.Recv()
				if err == io.EOF {
					log.Errorln("subscribe error eof.")
					goto Retry
				}
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Errorln("subscribe recv error")
					goto Retry
				}

				if resp != nil {
					pbStratumJob := resp.Job
					stratumJob := loadStratumJob(pbStratumJob)
					handler(stratumJob)
				} else {
					log.Errorln("subscribe resp is invalid")
					goto Retry
				}
			}
		} else {
			goto Retry
		}
	}
Out:
	log.Errorln("subscribe exit")
	return
}

func (s *GRPCService) SubmitShare(share *service.StratumShare, difficulty uint64) (service.ShareResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	var result service.ShareResult
	submitShareRequest := &pb.SubmitShareRequest{RegisterId: s.registerId, Share: loadPBStratumShare(share), Difficulty: difficulty}
	submitShareResponse, err := s.client.SubmitShare(ctx, submitShareRequest)

	if err != nil {
		result.State = service.StateUnknown
		return result, err
	}
	pbState := submitShareResponse.Result.State

	if pbState == pb.StratumShareState_SUC_SOLVE_SHARE || pbState == pb.StratumShareState_ERR_SUBMIT_BLOCK {
		result.State = service.StateSuccess
	} else if pbState == pb.StratumShareState_SUC_SUBMIT_BLOCK {
		result.State = service.StateSuccessSubmitBlock
	} else if pbState == pb.StratumShareState_ERR_DUPLICATE_SHARE {
		result.State = service.StateErrDuplicateShare
	} else if pbState == pb.StratumShareState_ERR_LOW_DIFFICULTY_SHARE {
		result.State = service.StateErrLowDifficultyShare
	} else {
		result.State = service.StateErrJobNotFound
	}

	if result.State == service.StateSuccess || result.State == service.StateSuccessSubmitBlock {
		result.ComputePower = service.ShareComputePower(submitShareResponse.Result.GetComputePower())
		if result.State == service.StateSuccessSubmitBlock {
			result.BlockHash = service.ShareBlockHash(submitShareResponse.Result.Hash)
		}

	}
	return result, nil
}

func (s *GRPCService) ClearShareHistory(height int32) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	clearShareHistoryRequest := &pb.ClearShareHistoryRequest{RegisterId: s.registerId, Height: height}
	clearShareHistoryResponse, err := s.client.ClearShareHistory(ctx, clearShareHistoryRequest)

	if err != nil {
		return false, err
	}

	return clearShareHistoryResponse.Result, nil
}

func (s *GRPCService) UnRegister() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	unRegisterRequest := &pb.UnRegisterRequest{RegisterId: s.registerId}
	unRegisterResponse, err := s.client.UnRegister(ctx, unRegisterRequest)

	if err != nil {
		return false, err
	}
	return unRegisterResponse.Result, nil
}

func loadStratumJob(pbJob *pb.StratumJob) *service.StratumJob {
	pbJobMeta := pbJob.Meta

	stratumJob := service.NewStratumJob()
	stratumJob.PrevHash = pbJob.PrevHash
	stratumJob.Version = pbJob.Version
	stratumJob.CoinBase1 = pbJob.CoinBase1
	stratumJob.CoinBase2 = pbJob.CoinBase2
	stratumJob.MerkleBranch = pbJob.MerkleBranch
	stratumJob.NBits = pbJob.NBits

	stratumJobMeta := new(service.StratumJobMeta)
	stratumJobMeta.Height = pbJobMeta.Height
	stratumJobMeta.CurTimeTs = pbJobMeta.CurTimeTs
	stratumJobMeta.MinTimeTs = pbJobMeta.MinTimeTs

	stratumJob.Meta = stratumJobMeta
	return stratumJob
}

func loadPBStratumShare(share *service.StratumShare) *pb.StratumShare {
	stratumShareMeta := share.Meta
	pbStratumShare := new(pb.StratumShare)
	pbStratumShare.NTime = share.NTime
	pbStratumShare.ExtraNonce2 = share.ExtraNonce2
	pbStratumShare.ExtraNonce1 = share.ExtraNonce1
	pbStratumShare.CoinBase1 = share.CoinBase1
	pbStratumShare.CoinBase2 = share.CoinBase2
	pbStratumShare.PrevHash = share.PrevHash
	pbStratumShare.Nonce = share.Nonce
	pbStratumShare.NTime = share.NTime
	pbStratumShare.Version = share.Version
	pbStratumShare.NBits = share.NBits

	pbStratumShareMeta := new(pb.StratumShareMeta)
	pbStratumShareMeta.Height = stratumShareMeta.Height
	pbStratumShareMeta.CurTimeTs = stratumShareMeta.JobCurTs
	pbStratumShare.Meta = pbStratumShareMeta

	return pbStratumShare
}
