package service

import (
	"fmt"
)

const (
	PrefixRegisterKey = "register_"
	PrefixJobKey      = "job_"
	PrefixShareKey    = "share_"
)

type CacheService interface {
	SetRegisterContext(key RegisterKey, r *Register) error
	GetRegisterContext(key RegisterKey) (*Register, error)
	DelRegisterContext(key RegisterKey) error
	SetBlockTransactions(key JobKey, expireTs int, transactions []*BlockTransactionPart) error
	GetBlockTransactions(key JobKey) ([]*BlockTransactionPart, error)
	SetShareHash(key ShareKey, hash string) error
	ExistShareHash(key ShareKey, hash string) (bool, error)
	ClearShareHistory(key ShareKey) error
}

type ShareKey string
type RegisterKey string
type JobKey string

func GenRegisterKey(registerId string) RegisterKey {
	return RegisterKey(fmt.Sprintf("%s_%s", PrefixRegisterKey, registerId))
}

func GenJobKey(registerId string, height int32, curTimeTs int32) JobKey {
	return JobKey(fmt.Sprintf("%s_%s_%d_%d", PrefixJobKey, registerId, height, curTimeTs))
}

func GenShareKey(registerId string, height int32) ShareKey {
	return ShareKey(fmt.Sprintf("%s_%s_%d", PrefixJobKey, registerId, height))
}
