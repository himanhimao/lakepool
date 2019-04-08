package service

import (
	"strconv"
	"os"
	"errors"
	"fmt"
	"encoding/hex"
	"strings"
	"context"
)

const (
	SeparatorJobId                        = "_"
	StateUnknown               ShareState = iota
	StateSuccess
	StateSuccessSubmitBlock
	StateErrDuplicateShare
	StateErrLowDifficultyShare
	StateErrJobNotFound
)

var (
	ErrorExtractJobId = errors.New("extract job id error.")
)

type StratumService interface {
	Init() error
	Register(config *StratumConfig, info *SysInfo) error
	GetLatestStratumJob() (*StratumJob, error)
	Subscribe(ctx context.Context, handler SubscribeHandler)
	SubmitShare(share *StratumShare, difficulty uint64) (ShareResult, error)
	ClearShareHistory(height int32) (bool, error)
	UnRegister() (bool, error)
}

type StratumConfig struct {
	coinType          string
	poolTag           string
	payoutAddress     string
	extraNonce1Length int
	extraNonce2Length int
}

type StratumJob struct {
	jobId        string
	prevHash     string
	coinBase1    string
	coinBase2    string
	merkleBranch []string
	version      string
	nBits        string
	nTime        string
	cleanJobs    bool
	meta         *StratumJobMeta
}

type StratumJobMeta struct {
	height    int32
	minTimeTs int32
	curTimeTs int32
}

type StratumShareHeader struct {
	version  string
	nTime    string
	nBits    string
	nonce    string
	prevHash string
}

type StratumShareCoinBase struct {
	extraNonce1 string
	extraNonce2 string
	coinBase1   string
	coinBase2   string
}

type StratumShare struct {
	StratumShareHeader
	StratumShareCoinBase
	meta *StratumShareMeta
}

type ShareState int
type ShareComputePower float64
type ShareBlockHash string

type ShareResult struct {
	State        ShareState
	ComputePower ShareComputePower
	BlockHash    ShareBlockHash
}

type StratumShareMeta struct {
	Height   int32
	JobCurTs int32
}

type SysInfo struct {
	pid      int
	hostName string
}

type SubscribeHandler func(job *StratumJob)

func NewSysInfo() *SysInfo {
	return &SysInfo{}
}

func NewStratumJob() *StratumJob {
	return &StratumJob{}
}

func NewStratumJobMeta() *StratumJobMeta {
	return &StratumJobMeta{}
}

func NewStratumShare() *StratumShare {
	return &StratumShare{}
}

func NewStratumShareMeta() *StratumShareMeta {
	return &StratumShareMeta{}
}

func (meta *StratumShareMeta) GetHeight() int32 {
	return meta.Height
}

func (meta *StratumShareMeta) GetJobCurTs() int32 {
	return meta.JobCurTs
}

func (ctx *SysInfo) GetPid() int {
	if ctx.pid == 0 {
		ctx.pid = os.Getpid()
	}
	return ctx.pid
}

func (ctx *SysInfo) GetHostName() (string, error) {
	if len(ctx.hostName) == 0 {
		if hostname, err := os.Hostname(); err != nil {
			return "", err
		} else {
			ctx.hostName = hostname
		}
	}
	return ctx.hostName, nil
}

func (job *StratumJob) Fill(jobId string, ts int64, cleanJobs bool) *StratumJob {
	job.cleanJobs = cleanJobs
	job.nTime = strconv.FormatInt(ts, 16)
	job.jobId = jobId
	return job
}

func (job *StratumJob) ToJSONInterface() ([]interface{}) {
	obj := make([]interface{}, 9)
	obj[0] = job.jobId
	obj[1] = job.prevHash
	obj[2] = job.coinBase1
	obj[3] = job.coinBase2
	obj[4] = job.merkleBranch
	obj[5] = job.version
	obj[6] = job.nBits
	obj[7] = job.nTime
	obj[8] = job.cleanJobs
	return obj

}

func (job *StratumJob) ToShare() *StratumShare {
	share := NewStratumShare().SetCoinBase1(job.GetCoinBase1()).SetCoinBase2(job.GetCoinBase2()).SetNBits(job.GetNBits()).
		SetPervHash(job.GetPrevHash()).SetVersion(job.GetVersion())
	meta := NewStratumShareMeta().SetJobCurTs(job.GetMeta().GetCurTimeTs()).SetHeight(job.GetMeta().GetHeight())
	share.SetMeta(meta)
	return share
}

func NewStratumConfig() *StratumConfig {
	return &StratumConfig{}
}

func (job *StratumJob) Resolve(num int) ([]*StratumJob, int) {
	return nil, 0
}

func (job *StratumJob) GetJobId() string {
	return job.jobId
}

func (job *StratumJob) GetPrevHash() string {
	return job.prevHash
}

func (job *StratumJob) GetCoinBase1() string {
	return job.coinBase1
}

func (job *StratumJob) GetCoinBase2() string {
	return job.coinBase2
}

func (job *StratumJob) GetMerkleBranch() []string {
	return job.merkleBranch
}

func (job *StratumJob) GetVersion() string {
	return job.version
}

func (job *StratumJob) GetNBits() string {
	return job.nBits
}

func (job *StratumJob) GetNTime() string {
	return job.nTime
}

func (job *StratumJob) GetMeta() *StratumJobMeta {
	return job.meta
}

func (job *StratumJob) IsCleanJobs() bool {
	return job.cleanJobs
}

func (meta *StratumJobMeta) GetHeight() int32 {
	return meta.height
}

func (meta *StratumJobMeta) GetMinTimeTs() int32 {
	return meta.minTimeTs
}

func (meta *StratumJobMeta) GetCurTimeTs() int32 {
	return meta.curTimeTs
}

func (job *StratumJob) SetPrevHash(prevHash string) *StratumJob {
	job.prevHash = prevHash
	return job
}

func (job *StratumJob) SetCoinBase1(coinbase1 string) *StratumJob {
	job.coinBase1 = coinbase1
	return job
}

func (job *StratumJob) SetCoinBase2(coinbase2 string) *StratumJob {
	job.coinBase2 = coinbase2
	return job
}

func (job *StratumJob) SetMerkleBranch(merkleBranch []string) *StratumJob {
	job.merkleBranch = merkleBranch
	return job
}

func (job *StratumJob) SetVersion(version string) *StratumJob {
	job.version = version
	return job
}

func (job *StratumJob) SetNBits(nBits string) *StratumJob {
	job.nBits = nBits
	return job
}

func (s *StratumShare) GetVersion() string {
	return s.version
}

func (s *StratumShare) GetExtraNonce1() string {
	return s.extraNonce1
}

func (s *StratumShare) GetExtraNonce2() string {
	return s.extraNonce2
}

func (s *StratumShare) GetNTime() string {
	return s.nTime
}

func (s *StratumShare) GetNBits() string {
	return s.nBits
}

func (s *StratumShare) GetNonce() string {
	return s.nonce
}

func (s *StratumShare) GetPervHash() string {
	return s.prevHash
}

func (s *StratumShare) GetCoinBase1() string {
	return s.coinBase1
}

func (s *StratumShare) GetCoinBase2() string {
	return s.coinBase2
}

func (s *StratumShare) GetMeta() *StratumShareMeta {
	return s.meta
}

func (job *StratumJob) SetMeta(meta *StratumJobMeta) *StratumJob {
	job.meta = meta
	return job
}

func (meta *StratumJobMeta) SetHeight(height int32) *StratumJobMeta {
	meta.height = height
	return meta
}

func (meta *StratumJobMeta) SetCurTimeTs(curTimeTs int32) *StratumJobMeta {
	meta.curTimeTs = curTimeTs
	return meta
}

func (meta *StratumJobMeta) SetMinTimeTs(minTimeTs int32) *StratumJobMeta {
	meta.minTimeTs = minTimeTs
	return meta
}

func (conf *StratumConfig) SetPoolTag(coinTag string) *StratumConfig {
	conf.poolTag = coinTag
	return conf
}

func (conf *StratumConfig) SetPayoutAddress(payoutAddress string) *StratumConfig {
	conf.payoutAddress = payoutAddress
	return conf
}

func (conf *StratumConfig) SetCoinType(coinType string) *StratumConfig {
	conf.coinType = coinType
	return conf
}

func (conf *StratumConfig) SetExtraNonce1Length(length int) *StratumConfig {
	conf.extraNonce1Length = length
	return conf
}

func (conf *StratumConfig) SetExtraNonce2Length(length int) *StratumConfig {
	conf.extraNonce2Length = length
	return conf
}

func (conf *StratumConfig) GetPoolTag() string {
	return conf.poolTag
}

func (conf *StratumConfig) GetPayoutAddress() string {
	return conf.payoutAddress
}

func (conf *StratumConfig) GetExtraNonce1Length() int {
	return conf.extraNonce1Length
}

func (conf *StratumConfig) GetExtraNonce2Length() int {
	return conf.extraNonce2Length
}

func (conf *StratumConfig) GetCoinType() string {
	return conf.coinType
}

func (s *StratumShare) SetVersion(version string) *StratumShare {
	s.version = version
	return s
}

func (s *StratumShare) SetExtraNonce1(extraNonce1 string) *StratumShare {
	s.extraNonce1 = extraNonce1
	return s
}

func (s *StratumShare) SetExtraNonce2(extraNonce2 string) *StratumShare {
	s.extraNonce2 = extraNonce2
	return s
}

func (s *StratumShare) SetNTime(nTime string) *StratumShare {
	s.nTime = nTime
	return s
}

func (s *StratumShare) SetNBits(nBits string) *StratumShare {
	s.nBits = nBits
	return s
}

func (s *StratumShare) SetNonce(nonce string) *StratumShare {
	s.nonce = nonce
	return s
}

func (s *StratumShare) SetPervHash(hash string) *StratumShare {
	s.prevHash = hash
	return s
}

func (s *StratumShare) SetCoinBase1(coinBase1 string) *StratumShare {
	s.coinBase1 = coinBase1
	return s
}

func (s *StratumShare) SetCoinBase2(coinBase2 string) *StratumShare {
	s.coinBase2 = coinBase2
	return s
}

func (s *StratumShare) SetMeta(meta *StratumShareMeta) *StratumShare {
	s.meta = meta
	return s
}

func (meta *StratumShareMeta) SetHeight(height int32) *StratumShareMeta {
	meta.Height = height
	return meta
}

func (meta *StratumShareMeta) SetJobCurTs(curTs int32) *StratumShareMeta {
	meta.JobCurTs = curTs
	return meta
}

func GenerateJobId(height int32, index int, difficulty uint64) string {
	jobIdStr := fmt.Sprintf("%d%s%d%s%d", height, SeparatorJobId, index, SeparatorJobId, difficulty)
	return hex.EncodeToString([]byte(jobIdStr))
}

func ExtractJobId(jobId string) (int32, int, uint64, error) {
	jobIdStr, err := hex.DecodeString(jobId)
	if err != nil {
		return 0, 0, 0, ErrorExtractJobId
	}
	data := strings.Split(string(jobIdStr), SeparatorJobId)

	if len(data) != 3 {
		return 0, 0, 0, ErrorExtractJobId
	}

	height, _ := strconv.Atoi(data[0])
	index, _ := strconv.Atoi(data[1])
	difficulty, _ := strconv.ParseInt(data[2], 10, 64)

	return int32(height), index, uint64(difficulty), nil
}
