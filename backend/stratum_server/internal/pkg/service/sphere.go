package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	SeparatorJobId            = "_"
	StateUnknown   ShareState = iota
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
	CoinType          string
	PoolTag           string
	PayoutAddress     string
	ExtraNonce1Length int
	ExtraNonce2Length int
}

type StratumJob struct {
	JobId        string
	PrevHash     string
	CoinBase1    string
	CoinBase2    string
	MerkleBranch []string
	Version      string
	NBits        string
	NTime        string
	CleanJobs    bool
	Meta         *StratumJobMeta
}

type StratumJobMeta struct {
	Height    int32
	MinTimeTs int32
	CurTimeTs int32
}

type StratumShareHeader struct {
	Version  string
	NTime    string
	NBits    string
	Nonce    string
	PrevHash string
}

type StratumShareCoinBase struct {
	ExtraNonce1 string
	ExtraNonce2 string
	CoinBase1   string
	CoinBase2   string
}

type StratumShare struct {
	StratumShareHeader
	StratumShareCoinBase
	Meta *StratumShareMeta
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
	job.CleanJobs = cleanJobs
	job.NTime = strconv.FormatInt(ts, 16)
	job.JobId = jobId
	return job
}

func (job *StratumJob) ToJSONInterface() ([]interface{}) {
	obj := make([]interface{}, 9)
	obj[0] = job.JobId
	obj[1] = job.PrevHash
	obj[2] = job.CoinBase1
	obj[3] = job.CoinBase2
	obj[4] = job.MerkleBranch
	obj[5] = job.Version
	obj[6] = job.NBits
	obj[7] = job.NTime
	obj[8] = job.CleanJobs
	return obj

}

func (job *StratumJob) ToShare() *StratumShare {
	share := NewStratumShare()
	share.CoinBase1 = job.CoinBase1
	share.CoinBase2 = job.CoinBase2
	share.NBits = job.NBits
	share.PrevHash = job.PrevHash
	share.Version =job.Version
	meta := NewStratumShareMeta()
	meta.JobCurTs = job.Meta.CurTimeTs
	meta.Height = job.Meta.Height
	share.Meta = meta
	return share
}

func NewStratumConfig() *StratumConfig {
	return &StratumConfig{}
}

func GenerateJobId(prefix string, height int32, index int, difficulty uint64) string {
	jobIdStr := fmt.Sprintf("%s%d%s%d%s%d", SeparatorJobId, height, SeparatorJobId, index, SeparatorJobId, difficulty)
	return fmt.Sprintf("%s%s", prefix, hex.EncodeToString([]byte(jobIdStr)))
}

func ExtractJobId(jobId string) (string, int32, int, uint64, error) {
	jobIdStr, err := hex.DecodeString(jobId)
	if err != nil {
		return "", 0, 0, 0, ErrorExtractJobId
	}
	data := strings.Split(string(jobIdStr), SeparatorJobId)

	if len(data) != 4{
		return "", 0, 0, 0, ErrorExtractJobId
	}

	prefix  := data[0]
	height, _ := strconv.Atoi(data[1])
	index, _ := strconv.Atoi(data[2])
	difficulty, _ := strconv.ParseInt(data[3], 10, 64)

	return prefix, int32(height), index, uint64(difficulty), nil
}
