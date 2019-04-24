package sphere

import (
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_sphere"
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
	"io"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/conf"
	"google.golang.org/grpc/balancer/roundrobin"
)

var (
	ReconnectWaitTime = time.Second * 15
	Timeout           = time.Second * 30
)

type GRPCService struct {
	config     *conf.GRPCConfig
	client     pb.SphereClient
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
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		return err
	}
	s.client = pb.NewSphereClient(conn)
	return nil
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
	pbStratumConfig.CoinType = config.GetCoinType()
	pbStratumConfig.PayoutAddress = config.GetPayoutAddress()
	pbStratumConfig.PoolTag = config.GetPoolTag()
	pbStratumConfig.ExtraNonce1Length = int32(config.GetExtraNonce1Length())
	pbStratumConfig.ExtraNonce2Length = int32(config.GetExtraNonce2Length())
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
	log.Infoln("Subscribing to gbt job")
	for {
		getLatestStratumJobRequest := &pb.GetLatestStratumJobRequest{RegisterId: s.registerId}
		stream, err := s.client.Subscribe(context.Background(), getLatestStratumJobRequest)
		if err != nil {
			log.Errorln("subscribe error ", err.Error())
			time.Sleep(ReconnectWaitTime)
		}
		if stream != nil {
			select {
			case <-stream.Context().Done():
				log.Errorln("sphere subscribe server interrupt. try reconnect..")
				break
			case <-ctx.Done():
				log.Debugln("subscribe routine cancel.")
				goto Out
			default:
				resp, err := stream.Recv()
				if err == io.EOF {
					log.Errorln("subscribe error eof.")
					break
				}
				if err != nil {
					log.Errorln(err.Error())
				}

				if resp != nil {
					pbStratumJob := resp.Job
					stratumJob := loadStratumJob(pbStratumJob)
					handler(stratumJob)
				}
			}
		}
	}
Out:
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
	stratumJob.SetPrevHash(pbJob.PrevHash)
	stratumJob.SetVersion(pbJob.Version)
	stratumJob.SetCoinBase1(pbJob.CoinBase1)
	stratumJob.SetCoinBase2(pbJob.CoinBase2)
	stratumJob.SetMerkleBranch(pbJob.MerkleBranch)
	stratumJob.SetNBits(pbJob.NBits)

	stratumJobMeta := new(service.StratumJobMeta)
	stratumJobMeta.SetHeight(pbJobMeta.Height)
	stratumJobMeta.SetCurTimeTs(pbJobMeta.CurTimeTs)
	stratumJobMeta.SetMinTimeTs(pbJobMeta.MinTimeTs)

	stratumJob.SetMeta(stratumJobMeta)
	return stratumJob
}

func loadPBStratumShare(share *service.StratumShare) *pb.StratumShare {
	stratumShareMeta := share.GetMeta()

	pbStratumShare := new(pb.StratumShare)
	pbStratumShare.NTime = share.GetNTime()
	pbStratumShare.ExtraNonce2 = share.GetExtraNonce2()
	pbStratumShare.ExtraNonce1 = share.GetExtraNonce1()
	pbStratumShare.CoinBase1 = share.GetCoinBase1()
	pbStratumShare.CoinBase2 = share.GetCoinBase2()
	pbStratumShare.PrevHash = share.GetPervHash()
	pbStratumShare.Nonce = share.GetNonce()
	pbStratumShare.NTime = share.GetNTime()
	pbStratumShare.Version = share.GetVersion()
	pbStratumShare.NBits = share.GetNBits()

	pbStratumShareMeta := new(pb.StratumShareMeta)
	pbStratumShareMeta.Height = stratumShareMeta.GetHeight()
	pbStratumShareMeta.CurTimeTs = stratumShareMeta.GetJobCurTs()
	pbStratumShare.Meta = pbStratumShareMeta

	return pbStratumShare
}
