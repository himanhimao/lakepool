package service

const (
	KeyRegId                = "id"
	KeyRegPayoutAddress     = "payoutAddress"
	KeyRegPoolTag           = "poolTag"
	KeyRegCoinType          = "coinType"
	KeyRegUsedTestNet       = "usedTestNet"
	KeyRegExtraNonce1Length = "extraNonce1Length"
	KeyRegExtraNonce2Length = "extraNonce2Length"
)

type Register struct {
	Id                string
	PayoutAddress     string
	PoolTag           string
	CoinType          string
	UsedTestNet       bool
	ExtraNonce1Length int
	ExtraNonce2Length int
}


func (r *Register) IsValid() bool {
	return len(r.Id) > 0 && len(r.PayoutAddress) > 0 && len(r.CoinType) > 0 && len(r.PoolTag) > 0
}

func NewRegister() *Register {
	return &Register{}
}
