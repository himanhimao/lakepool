package service

import "math/big"

type CoinService interface {
	IsValidAddress(address string, isUsedTestNet bool) bool
	GetLatestStratumJob(registerId string, r *Register) (*StratumJobPart, []*BlockTransactionPart, error)
	MakeBlock(header *BlockHeaderPart, base *BlockCoinBasePart, transactions []*BlockTransactionPart) (*Block, error)
	SubmitBlock(data string) (bool, error)
	IsSolveHash(hash string, targetDifficulty *big.Int) (bool, error)
	GetTargetDifficulty(bits string) (*big.Int, error)
	CalculateShareComputePower(difficulty *big.Int) (*big.Int, error)
}

type Block struct {
	hash string
	data string
}

type StratumJobMetaPart struct {
	height    int32
	minTimeTs int32
	curTimeTs int32
}

type StratumJobPart struct {
	prevHash     string
	coinBase1    string
	coinBase2    string
	merkleBranch []string
	version      string
	nBits        string
	meta         *StratumJobMetaPart
}

type BlockTransactionPart struct {
	Data string
}

type BlockHeaderPart struct {
	version  string
	prevHash string
	nBits    string
	nonce    string
	nTime    string
}

type BlockCoinBasePart struct {
	coinBase1   string
	coinBase2   string
	extraNonce1 string
	extraNonce2 string
}

func NewStratumJobMetaPart() *StratumJobMetaPart {
	return &StratumJobMetaPart{}
}

func NewStratumJobPart() *StratumJobPart {
	return &StratumJobPart{}
}

func NewBlockTransactionPart() *BlockTransactionPart {
	return &BlockTransactionPart{}
}

func NewBlockCoinBasePart() *BlockCoinBasePart {
	return &BlockCoinBasePart{}
}

func NewBlockHeaderPart() *BlockHeaderPart {
	return &BlockHeaderPart{}
}

func NewBlock() *Block {
	return &Block{}
}

func (b *Block) SetHash(hash string) *Block {
	b.hash = hash
	return b
}

func (b *Block) SetData(data string) *Block {
	b.data = data
	return b
}

func (b *Block) GetHash() string {
	return b.hash
}

func (b *Block) GetData() string {
	return b.data
}

func (s *StratumJobPart) SetPrevHash(prevHash string) *StratumJobPart {
	s.prevHash = prevHash
	return s
}

func (s *StratumJobPart) SetCoinBase1(coinBase1 string) *StratumJobPart {
	s.coinBase1 = coinBase1
	return s
}

func (s *StratumJobPart) SetCoinBase2(coinbase2 string) *StratumJobPart {
	s.coinBase2 = coinbase2
	return s
}

func (s *StratumJobPart) SetMerkleBranch(merkleBranch []string) *StratumJobPart {
	s.merkleBranch = merkleBranch
	return s
}

func (s *StratumJobPart) SetVersion(version string) *StratumJobPart {
	s.version = version
	return s
}

func (s *StratumJobPart) SetNBits(nbits string) *StratumJobPart {
	s.nBits = nbits
	return s
}

func (s *StratumJobPart) SetMeta(meta *StratumJobMetaPart) *StratumJobPart {
	s.meta = meta
	return s
}

func (m *StratumJobMetaPart) SetHeight(Height int32) *StratumJobMetaPart {
	m.height = Height
	return m
}

func (m *StratumJobMetaPart) SetMinTimeTs(minTimeTs int32) *StratumJobMetaPart {
	m.minTimeTs = minTimeTs
	return m
}

func (m *StratumJobMetaPart) SetCurTimeTs(curTimeTs int32) *StratumJobMetaPart {
	m.curTimeTs = curTimeTs
	return m
}

func (t *BlockTransactionPart) GetData() string {
	return t.Data
}

func (t *BlockTransactionPart) SetData(data string) *BlockTransactionPart {
	t.Data = data
	return t
}

func (s *StratumJobPart) GetPrevHash() string {
	return s.prevHash
}

func (s *StratumJobPart) GetCoinBase1() string {
	return s.coinBase1
}

func (s *StratumJobPart) GetCoinBase2() string {
	return s.coinBase2
}

func (s *StratumJobPart) GetMerkleBranch() []string {
	return s.merkleBranch
}

func (s *StratumJobPart) GetVersion() string {
	return s.version
}

func (s *StratumJobPart) GetNBits() string {
	return s.nBits
}

func (s *StratumJobPart) GetMeta() *StratumJobMetaPart {
	return s.meta
}

func (m *StratumJobMetaPart) GetHeight() int32 {
	return m.height
}

func (m *StratumJobMetaPart) GetMinTimeTs() int32 {
	return m.minTimeTs
}

func (m *StratumJobMetaPart) GetCurTimeTs() int32 {
	return m.curTimeTs
}

func (b *BlockCoinBasePart) GetCoinBase1() string {
	return b.coinBase1
}

func (b *BlockCoinBasePart) GetCoinBase2() string {
	return b.coinBase2
}

func (b *BlockCoinBasePart) GetExtraNonce1() string {
	return b.extraNonce1
}

func (b *BlockCoinBasePart) GetExtraNonce2() string {
	return b.extraNonce2
}

func (b *BlockCoinBasePart) SetCoinBase1(coinBase1 string) *BlockCoinBasePart {
	b.coinBase1 = coinBase1
	return b
}

func (b *BlockCoinBasePart) SetCoinBase2(coinBase2 string) *BlockCoinBasePart {
	b.coinBase2 = coinBase2
	return b
}

func (b *BlockCoinBasePart) SetExtraNonce1(extraNonce1 string) *BlockCoinBasePart {
	b.extraNonce1 = extraNonce1
	return b
}

func (b *BlockCoinBasePart) SetExtraNonce2(extraNonce2 string) *BlockCoinBasePart {
	b.extraNonce2 = extraNonce2
	return b
}

func (b *BlockHeaderPart) SetVersion(version string) *BlockHeaderPart {
	b.version = version
	return b
}

func (b *BlockHeaderPart) SetPrevHash(prevHash string) *BlockHeaderPart {
	b.prevHash = prevHash
	return b
}

func (b *BlockHeaderPart) SetNBits(nBits string) *BlockHeaderPart {
	b.nBits = nBits
	return b
}

func (b *BlockHeaderPart) SetNonce(nonce string) *BlockHeaderPart {
	b.nonce = nonce
	return b
}

func (b *BlockHeaderPart) SetNTime(nTime string) *BlockHeaderPart {
	b.nTime = nTime
	return b
}

func (b *BlockHeaderPart) GetVersion() string {
	return b.version
}

func (b *BlockHeaderPart) GetPrevHash() string {
	return b.prevHash
}

func (b *BlockHeaderPart) GetNBits() string {
	return b.nBits
}

func (b *BlockHeaderPart) GetNonce() string {
	return b.nonce
}

func (b *BlockHeaderPart) GetNTime() string {
	return b.nTime
}
