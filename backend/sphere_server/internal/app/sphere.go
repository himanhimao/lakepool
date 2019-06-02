package impl

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/conf"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/service"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_sphere"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"strconv"
	"sync"
	"time"
)

const (
	maxPid       = 4194303
	requestIdLen = 8
)

type SphereServer struct {
	Conf       *conf.SphereConfig
	Mgr        *service.Manager
	contextMap sync.Map
}

func (s *SphereServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if len(in.SysInfo.Hostname) == 0 {
		st := status.New(codes.InvalidArgument, "Invalid argument - hostname")
		return nil, st.Err()
	}

	if in.SysInfo.Pid < 0 || in.SysInfo.Pid > maxPid {
		st := status.New(codes.InvalidArgument, fmt.Sprintf("%s:%d", "Invalid argument - pid", in.SysInfo.Pid))
		return nil, st.Err()
	}

	var coinService service.CoinService
	if coinService = s.Mgr.GetCoinService(in.Config.CoinType); coinService == nil {
		st := status.New(codes.InvalidArgument, fmt.Sprintf("%s:%s", "Invalid argument - coinType", in.Config.CoinType))
		return nil, st.Err()
	}

	if len(in.Config.PayoutAddress) == 0 || !coinService.IsValidAddress(in.Config.PayoutAddress, in.Config.UsedTestNet) {
		st := status.New(codes.InvalidArgument, "Invalid argument - payoutAddress")
		return nil, st.Err()
	}

	if len(in.Config.PoolTag) == 0 {
		st := status.New(codes.InvalidArgument, "Invalid argument - poolTag")
		return nil, st.Err()
	}

	registerId := s.calculateRegisterId(in.SysInfo)
	registerKey := service.GenRegisterKey(registerId)
	register := service.NewRegister()
	register.PayoutAddress = in.Config.PayoutAddress
	register.PoolTag = in.Config.PoolTag
	register.CoinType = in.Config.CoinType
	register.UsedTestNet = in.Config.UsedTestNet
	register.ExtraNonce1Length = int(in.Config.ExtraNonce1Length)
	register.ExtraNonce2Length = int(in.Config.ExtraNonce2Length)

	if err := s.storeRegisterContext(registerKey, register); err != nil {
		st := status.New(codes.Internal, fmt.Sprintf("Store context error: %s", err.Error()))
		return nil, st.Err()
	}
	return &pb.RegisterResponse{RegisterId: registerId}, nil
}

func (s *SphereServer) GetLatestStratumJob(ctx context.Context, in *pb.GetLatestStratumJobRequest) (*pb.GetLatestStratumJobResponse, error) {
	registerId := in.RegisterId
	if len(registerId) != requestIdLen {
		st := status.New(codes.InvalidArgument, "Abnormal - invalid argument requestId")
		return nil, st.Err()
	}
	var err error
	var register *service.Register
	var coinService service.CoinService
	var cacheService service.CacheService

	registerKey := service.GenRegisterKey(registerId)
	register, err = s.fetchRegisterContext(registerKey)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if register == nil || !register.IsValid() {
		st := status.New(codes.InvalidArgument, "Abnormal - unknown requestId")
		return nil, st.Err()
	}

	coinService = s.Mgr.GetCoinService(register.CoinType)
	if coinService == nil {
		st := status.New(codes.Internal, "Abnormal - coin service")
		return nil, st.Err()
	}

	stratumJobPart, jobTransactions, err := coinService.GetLatestStratumJob(registerId, register)
	if err != nil {
		st := status.New(codes.Internal, fmt.Sprintf("%s : %s", "Abnormal - get stratum job ", err))
		return nil, st.Err()
	}

	cacheService = s.Mgr.GetCacheService()
	jobKey := service.GenJobKey(registerId, stratumJobPart.Meta.Height, stratumJobPart.Meta.CurTimeTs)
	jobCacheExpireTs := s.Conf.Configs[register.CoinType].JobCacheExpireTs.Seconds()
	if cacheService == nil {
		st := status.New(codes.Internal, "Abnormal - cache service")
		return nil, st.Err()
	}

	err = cacheService.SetBlockTransactions(jobKey, int(jobCacheExpireTs), jobTransactions)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	return &pb.GetLatestStratumJobResponse{Job: stratumJobPart.ToPBStratumJob()}, nil
}

func (s *SphereServer) SubmitShare(ctx context.Context, in *pb.SubmitShareRequest) (*pb.SubmitShareResponse, error) {
	registerId := in.RegisterId
	if len(registerId) != requestIdLen {
		st := status.New(codes.InvalidArgument, "Abnormal - invalid argument requestId")
		return nil, st.Err()
	}
	var err error
	var register *service.Register
	var coinService service.CoinService
	var cacheService service.CacheService

	registerKey := service.GenRegisterKey(registerId)
	register, err = s.fetchRegisterContext(registerKey)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if register == nil || !register.IsValid() {
		st := status.New(codes.InvalidArgument, "Abnormal - unknown requestId")
		return nil, st.Err()
	}

	coinService = s.Mgr.GetCoinService(register.CoinType)
	stratumShare := in.Share
	jobKey := service.GenJobKey(registerId, stratumShare.Meta.Height, stratumShare.Meta.CurTimeTs)
	if coinService == nil {
		st := status.New(codes.Internal, "Abnormal - coin service")
		return nil, st.Err()
	}

	cacheService = s.Mgr.GetCacheService()
	if cacheService == nil {
		st := status.New(codes.Internal, "Abnormal - cache service")
		return nil, st.Err()
	}

	transactions, err := cacheService.GetBlockTransactions(jobKey)
	if err != nil {
		st := status.New(codes.Internal, fmt.Sprintf("Abnormal - get block transactions:%s", err.Error()))
		return nil, st.Err()
	}

	if transactions == nil {
		return &pb.SubmitShareResponse{Result: &pb.SubmitShareResult{State: pb.StratumShareState_ERR_JOB_NOT_FOUND}}, nil
	}

	blockHeaderPart, coinBasePart := splitServiceParts(in.Share)
	block, err := coinService.MakeBlock(blockHeaderPart, coinBasePart, transactions)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	targetDifficulty := new(big.Int).SetUint64(in.Difficulty)
	isSolveHash, err := coinService.IsSolveHash(block.Hash, targetDifficulty)

	if err != nil {
		st := status.New(codes.Internal, fmt.Sprintf("Abnormal - is slove hash:%s", err.Error()))
		return nil, st.Err()
	}

	if !isSolveHash {
		return &pb.SubmitShareResponse{Result: &pb.SubmitShareResult{State: pb.StratumShareState_ERR_LOW_DIFFICULTY_SHARE}}, nil
	}

	//duplicate Check
	shareKey := service.GenShareKey(registerId, in.Share.Meta.Height)
	isExistShare, err := cacheService.ExistShareHash(shareKey, block.Hash)
	if err != nil {
		st := status.New(codes.Internal, fmt.Sprintf("Abnormal - is exist share:%s", err.Error()))
		return nil, st.Err()
	}

	if isExistShare {
		return &pb.SubmitShareResponse{Result: &pb.SubmitShareResult{State: pb.StratumShareState_ERR_DUPLICATE_SHARE}}, nil
	}

	netTargetDifficulty, err := coinService.GetTargetDifficulty(blockHeaderPart.NBits)
	if err != nil {
		st := status.New(codes.Internal, fmt.Sprintf("Abnormal - get target difficulty:%s", err.Error()))
		return nil, st.Err()
	}

	isSubmitHash, err := coinService.IsSolveHash(block.Hash, netTargetDifficulty)
	if err != nil {
		st := status.New(codes.Internal, fmt.Sprintf("Abnormal - is submit hash:%s", err.Error()))
		return nil, st.Err()
	}

	var state pb.StratumShareState
	var submitState bool
	if isSubmitHash {
		submitState, err = coinService.SubmitBlock(block.Data)
		if err != nil {
			st := status.New(codes.Internal, fmt.Sprintf("Abnormal - submit share :%s", err.Error()))
			return nil, st.Err()
		}
	}

	if submitState {
		state = pb.StratumShareState_SUC_SUBMIT_BLOCK
	} else {
		state = pb.StratumShareState_ERR_SUBMIT_BLOCK
	}
	shareComputePower, err := coinService.CalculateShareComputePower(targetDifficulty)
	if err != nil {
		st := status.New(codes.Internal, fmt.Sprintf("Abnormal - calculate share compute:%s", err.Error()))
		return nil, st.Err()
	}

	return &pb.SubmitShareResponse{Result: &pb.SubmitShareResult{State: state, Hash: block.Hash, ComputePower: float64(shareComputePower.Uint64())}}, nil
}

func (s *SphereServer) Subscribe(in *pb.GetLatestStratumJobRequest, stream pb.Sphere_SubscribeServer) error {
	registerId := in.RegisterId
	if len(registerId) != requestIdLen {
		st := status.New(codes.InvalidArgument, "Abnormal - invalid argument requestId")
		return st.Err()
	}
	var err error
	var register *service.Register
	var coinService service.CoinService
	var cacheService service.CacheService

	registerKey := service.GenRegisterKey(registerId)
	register, err = s.fetchRegisterContext(registerKey)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return st.Err()
	}

	if register == nil || !register.IsValid() {
		st := status.New(codes.InvalidArgument, "Abnormal - unknown requestId")
		return st.Err()
	}

	coinService = s.Mgr.GetCoinService(register.CoinType)
	if coinService == nil {
		st := status.New(codes.Internal, "Abnormal - coin service")
		return st.Err()
	}

	cacheService = s.Mgr.GetCacheService()
	if cacheService == nil {
		st := status.New(codes.Internal, "Abnormal - cache service")
		return st.Err()
	}

	pullGBTInterval := s.Conf.Configs[register.CoinType].PullGBTInterval
	notifyInterval := s.Conf.Configs[register.CoinType].NotifyInterval
	notifyTicker := time.NewTicker(notifyInterval)
	pullGBTTicker := time.NewTicker(pullGBTInterval)
	var currentBlockHeight int
	var stratumJobPart *service.StratumJobPart
	var blockTransactions []*service.BlockTransactionPart
	notifyChan := make(chan bool, 1)
	jobCacheExpireTs := int(s.Conf.Configs[register.CoinType].JobCacheExpireTs.Seconds())

	log.Infoln("pull GBT interval: ", pullGBTInterval.Seconds(), "notify interval: ", notifyInterval.Seconds())
	for {
		select {
		case <-stream.Context().Done():
			log.Debugln("stream context done")
			goto Out
		case <-notifyTicker.C:
			notifyChan <- true
		case <-notifyChan:
			stratumJobPart, blockTransactions, err = coinService.GetLatestStratumJob(registerId, register)
			if err != nil {
				log.WithFields(log.Fields{
					"error":       err,
					"register_id": registerId,
					"coinType":    register.CoinType,
				}).Errorln("get latest stratum job error")
				break
			}

			jobKey := service.GenJobKey(registerId, stratumJobPart.Meta.Height, stratumJobPart.Meta.CurTimeTs)
			err = cacheService.SetBlockTransactions(jobKey, jobCacheExpireTs, blockTransactions)
			if err != nil {
				log.WithFields(log.Fields{
					"error":       err,
					"register_id": registerId,
					"coinType":    register.CoinType,
				}).Errorln("set block transactions error")
				break
			}
			stream.Send(&pb.GetLatestStratumJobResponse{Job: stratumJobPart.ToPBStratumJob()})
			log.WithFields(log.Fields{
				"block_height": stratumJobPart.Meta.Height,
				"register_id":  registerId,
				"coinType":     register.CoinType,
			}).Infoln("job has been sent")
		case <-pullGBTTicker.C:
			newBlockHeight, err := coinService.GetNewBlockHeight()
			if err != nil {
				log.WithFields(log.Fields{
					"error":       err,
					"register_id": registerId,
					"coinType":    register.CoinType,
				}).Errorln("get new block height  error")
				break
			}
			if newBlockHeight != currentBlockHeight {
				log.WithFields(log.Fields{
					"new_block_height": newBlockHeight,
					"cur_block_height": currentBlockHeight,
					"register_id":      registerId,
					"coinType":         register.CoinType,
				}).Infoln("new block height")
				currentBlockHeight = newBlockHeight
				notifyChan <- true
			}
		}
	}
Out:
	return nil
}

func (s *SphereServer) ClearShareHistory(ctx context.Context, in *pb.ClearShareHistoryRequest) (*pb.ClearShareHistoryResponse, error) {
	registerId := in.RegisterId
	if len(registerId) != requestIdLen {
		st := status.New(codes.InvalidArgument, "Abnormal - invalid argument requestId")
		return nil, st.Err()
	}
	var err error
	var register *service.Register
	var cacheService service.CacheService

	registerKey := service.GenRegisterKey(registerId)
	register, err = s.fetchRegisterContext(registerKey)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if register == nil || !register.IsValid() {
		st := status.New(codes.InvalidArgument, "Abnormal - unknown requestId")
		return nil, st.Err()
	}

	cacheService = s.Mgr.GetCacheService()
	if cacheService == nil {
		st := status.New(codes.Internal, "Abnormal - cache service")
		return nil, st.Err()
	}

	err = cacheService.ClearShareHistory(service.GenShareKey(registerId, in.Height))
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	return &pb.ClearShareHistoryResponse{Result: true}, nil
}

func (s *SphereServer) UnRegister(ctx context.Context, in *pb.UnRegisterRequest, ) (*pb.UnRegisterResponse, error) {
	registerId := in.RegisterId
	if len(registerId) != requestIdLen {
		st := status.New(codes.InvalidArgument, "Abnormal - Invalid argument requestId")
		return nil, st.Err()
	}
	var err error
	var register *service.Register
	var cacheService service.CacheService

	registerKey := service.GenRegisterKey(registerId)
	register, err = s.fetchRegisterContext(registerKey)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if register == nil || !register.IsValid() {
		st := status.New(codes.InvalidArgument, "Abnormal - unknown requestId")
		return nil, st.Err()
	}

	cacheService = s.Mgr.GetCacheService()
	if cacheService == nil {
		st := status.New(codes.Internal, "Abnormal - cache service")
		return nil, st.Err()
	}

	err = cacheService.DelRegisterContext(service.GenRegisterKey(registerId))
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}
	s.contextMap.Store(registerKey, nil)
	return &pb.UnRegisterResponse{Result: true}, nil
}

func (s *SphereServer) calculateRegisterId(info *pb.SysInfo) string {
	buf := new(bytes.Buffer)
	buf.WriteString(info.Hostname)
	buf.WriteString(strconv.Itoa(int(info.Pid)))
	hash := md5.Sum(buf.Bytes())
	registerId := fmt.Sprintf("%x", hash)
	return registerId[3:11]
}

func (s *SphereServer) storeRegisterContext(registerKey service.RegisterKey, r *service.Register) error {
	s.contextMap.Store(registerKey, r)
	return s.Mgr.GetCacheService().SetRegisterContext(registerKey, r)
}

func (s *SphereServer) fetchRegisterContext(registerKey service.RegisterKey) (*service.Register, error) {
	if value, ok := s.contextMap.Load(registerKey); ok {
		return value.(*service.Register), nil
	}
	return s.Mgr.GetCacheService().GetRegisterContext(registerKey)
}

func splitServiceParts(share *pb.StratumShare) (*service.BlockHeaderPart, *service.BlockCoinBasePart) {
	coinBasePart := service.NewBlockCoinBasePart()
	coinBasePart.CoinBase1 = share.CoinBase1
	coinBasePart.CoinBase2 = share.CoinBase2
	coinBasePart.ExtraNonce1 = share.ExtraNonce1
	coinBasePart.ExtraNonce2 = share.ExtraNonce2

	blockHeaderPart := service.NewBlockHeaderPart()
	blockHeaderPart.NTime = share.NTime
	blockHeaderPart.NBits = share.NBits
	blockHeaderPart.Version = share.Version
	blockHeaderPart.Nonce = share.Nonce
	blockHeaderPart.PrevHash = share.PrevHash
	return blockHeaderPart, coinBasePart
}
