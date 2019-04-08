package service

const (
	KeyRegPayoutAddress     = "payoutAddress"
	KeyRegPoolTag           = "poolTag"
	KeyRegCoinType          = "coinType"
	KeyRegUsedTestNet       = "usedTestNet"
	KeyRegExtraNonce1Length = "extraNonce1Length"
	KeyRegExtraNonce2Length = "extraNonce2Length"
)

type Register struct {
	payoutAddress     string
	poolTag           string
	coinType          string
	usedTestNet       bool
	extraNonce1Length int
	extraNonce2Length int
}

func (r *Register) GetPayoutAddress() string {
	return r.payoutAddress
}

func (r *Register) SetPayoutAddress(coinAddress string) *Register {
	r.payoutAddress = coinAddress
	return r
}

func (r *Register) GetPoolTag() string {
	return r.poolTag
}

func (r *Register) SetPoolTag(poolTag string) *Register {
	r.poolTag = poolTag
	return r
}

func (r *Register) GetCoinType() string {
	return r.coinType
}

func (r *Register) SetCoinType(coinType string) *Register {
	r.coinType = coinType
	return r
}

func (r *Register) IsUsedTestNet() bool {
	return r.usedTestNet
}

func (r *Register) SetUsedTestNet(result bool) *Register {
	r.usedTestNet = result
	return r
}

func (r *Register) IsValid() bool {
	return len(r.payoutAddress) > 0 && len(r.coinType) > 0 && len(r.poolTag) > 0
}

func (r *Register) GetExtraNonce1Length() int {
	return r.extraNonce1Length
}

func (r *Register) SetExtraNonce1Length(length int) *Register {
	r.extraNonce1Length = length
	return r
}

func (r *Register) GetExtraNonce2Length() int {
	return r.extraNonce2Length
}

func (r *Register) SetExtraNonce2Length(length int) *Register {
	r.extraNonce2Length = length
	return r
}

func NewRegister() *Register {
	return &Register{}
}
