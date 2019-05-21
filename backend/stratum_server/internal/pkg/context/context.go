package context

import (
	"context"
	"sync"
	"github.com/prep/average"
)

type ConnContext struct {
	RemoteIP   string
	LocalIP    string
	AcceptedTs int64
	ClosedTs   int64
}

type StratumSubscribeContext struct {
	SessionID   string
	UserAgent   string
	SubscribeTs int64
}

type StratumAuthorizeContext struct {
	WorkerName  string
	ExtName     string
	UserName    string
	Password    string
	AuthorizeTs int64
}

type StratumSetDifficultyContext struct {
	LatestDifficulty uint64
}

type StratumNotifyContext struct {
	StartNotifyTs         int64
	LatestNotifyJobHeight int32
	LatestNotifyJobIndex  int
	LatestNotifyTs        int64
	NotifyCount           int64
	NotifyMutex           sync.Mutex
}

type StratumSubmitContext struct {
	SlideWindow    *average.SlidingWindow
	SubmitCount    int64
	AcceptCount    int64
	RejectCount    int64
	ErrorCount     int64
	LatestSubmitTs int64
}

type StratumContext struct {
	Context     context.Context
	CancelFunc  context.CancelFunc
	ConnContext
	StratumSubscribeContext
	StratumAuthorizeContext
	StratumSetDifficultyContext
	StratumNotifyContext
	StratumSubmitContext
}

func NewStratumContext() *StratumContext {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &StratumContext{Context: ctx, CancelFunc: cancel}
}
