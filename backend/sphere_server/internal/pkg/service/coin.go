package service

import (
	pb "github.com/himanhimao/lakepool_proto/backend/proto_sphere"
	"math/big"
)

type CoinService interface {
	IsValidAddress(address string, isUsedTestNet bool) bool
	GetLatestStratumJob(r *Register) (*StratumJobPart, []*Transaction, error)
	MakeBlock(r *Register, header *BlockHeaderPart, base *BlockCoinBasePart, transactions []*Transaction) (*Block, error)
	SubmitBlock(data string) (bool, error)
	IsSolveHash(hash string, targetDifficulty *big.Int) (bool, error)
	GetTargetDifficulty(bits string) (*big.Int, error)
	CalculateShareComputePower(difficulty *big.Int) (*big.Int, error)
	GetNewBlockHeight() (int, error)
}

type Block struct {
	Hash string
	Data string
}

type StratumJobMetaPart struct {
	Height    int32
	MinTimeTs int32
	CurTimeTs int32
}

type StratumJobPart struct {
	PrevHash     string
	CoinBase1    string
	CoinBase2    string
	MerkleBranch []string
	Version      string
	NBits        string
	Meta         *StratumJobMetaPart
}

type Transaction struct {
	Data string
	Hash string
}

type BlockHeaderPart struct {
	Version  string
	PrevHash string
	NBits    string
	Nonce    string
	NTime    string
}

type BlockCoinBasePart struct {
	CoinBase1   string
	CoinBase2   string
	ExtraNonce1 string
	ExtraNonce2 string
}

func NewStratumJobMetaPart() *StratumJobMetaPart {
	return &StratumJobMetaPart{}
}

func NewStratumJobPart() *StratumJobPart {
	return &StratumJobPart{}
}

func NewBlockTransactionPart(hash string, data string) *Transaction {
	return &Transaction{Hash: hash, Data: data}
}

func NewBlockCoinBasePart() *BlockCoinBasePart {
	return &BlockCoinBasePart{}
}

func NewBlockHeaderPart() *BlockHeaderPart {
	return &BlockHeaderPart{}
}

func NewBlock(hash string, data string) *Block {
	return &Block{Hash:hash,Data:data}
}

func (job *StratumJobPart) ToPBStratumJob() *pb.StratumJob {
	pbStratumJob := new(pb.StratumJob)
	pbStratumJob.NBits = job.NBits
	pbStratumJob.PrevHash = job.PrevHash
	pbStratumJob.MerkleBranch = job.MerkleBranch
	pbStratumJob.CoinBase1 = job.CoinBase1
	pbStratumJob.CoinBase2 = job.CoinBase2
	pbStratumJob.Version = job.Version

	pbStratumJobMeta := new(pb.StratumJobMeta)
	pbStratumJobMeta.CurTimeTs = job.Meta.CurTimeTs
	pbStratumJobMeta.Height = job.Meta.Height
	pbStratumJobMeta.MinTimeTs = job.Meta.MinTimeTs
	pbStratumJob.Meta = pbStratumJobMeta
	return pbStratumJob
}
