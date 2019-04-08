package context

import (
	"context"
	"sync"
	"github.com/prep/average"
)

type ConnContext struct {
	remoteIP   string
	localIP    string
	acceptedTs int64
	closedTs   int64
}

type StratumSubscribeContext struct {
	sessionID   string
	userAgent   string
	subscribeTs int64
}

type StratumAuthorizeContext struct {
	workerName  string
	extName     string
	userName    string
	password    string
	authorizeTs int64
}

type StratumSetDifficultyContext struct {
	latestDifficulty uint64
}

type StratumNotifyContext struct {
	startNotifyTs         int64
	latestNotifyJobHeight int32
	latestNotifyJobIndex  int
	latestNotifyTs        int64
	notifyCount           int64
	notifyMutex           sync.Mutex
}

type StratumSubmitContext struct {
	slideWindow    *average.SlidingWindow
	submitCount    int64
	acceptCount    int64
	rejectCount    int64
	errorCount     int64
	latestSubmitTs int64
}

type StratumContext struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
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
	return &StratumContext{ctx: ctx, cancelFunc: cancel}
}

func (ctx *StratumContext) CancelContext() {
	ctx.cancelFunc()
}

func (ctx *StratumContext) GetStartNotifyTs() int64 {
	return ctx.startNotifyTs
}

func (ctx *StratumContext) GetLatestNotifyJobHeight() int32 {
	return ctx.latestNotifyJobHeight
}

func (ctx *StratumContext) GetLatestNotifyJobIndex() int {
	return ctx.latestNotifyJobIndex
}

func (ctx *StratumContext) GetLatestNotifyTs() int64 {
	return ctx.latestNotifyTs
}

func (ctx *StratumContext) GetLatestDifficulty() uint64 {
	return ctx.latestDifficulty
}

func (ctx *StratumContext) GetAuthorizeWorkerName() string {
	return ctx.workerName
}

func (ctx *StratumContext) GetAuthorizeExtName() string {
	return ctx.extName
}

func (ctx *StratumContext) GetAuthorizeUserName() string {
	return ctx.userName
}

func (ctx *StratumContext) GetAuthorizePassword() string {
	return ctx.password
}

func (ctx *StratumContext) GetNotifyCount() int64 {
	return ctx.notifyCount
}

func (ctx *StratumContext) GetAuthorizeTs() int64 {
	return ctx.authorizeTs
}

func (ctx *StratumContext) GetSubscribeSessionID() string {
	return ctx.sessionID
}

func (ctx *StratumContext) GetSubscribeUserAgent() string {
	return ctx.userAgent
}

func (ctx *StratumContext) GetSubScribeTs() int64 {
	return ctx.subscribeTs
}

func (ctx *StratumContext) GetConnRemoteIP() string {
	return ctx.remoteIP
}

func (ctx *StratumContext) GetConnLocalIP() string {
	return ctx.localIP
}

func (ctx *StratumContext) GetConnAcceptedTs() int64 {
	return ctx.acceptedTs
}

func (ctx *StratumContext) GetConnClosedTs() int64 {
	return ctx.closedTs
}

func (ctx *StratumContext) GetSlideWindow() *average.SlidingWindow {
	return ctx.slideWindow
}

func (ctx *StratumContext) GetLatestSubmitTs() int64 {
	return ctx.latestSubmitTs
}

func (ctx *StratumContext) GetAcceptCount() int64 {
	return ctx.acceptCount
}

func (ctx *StratumContext) GetSubmitCount() int64 {
	return ctx.acceptCount
}

func (ctx *StratumContext) GetRejectCount() int64 {
	return ctx.rejectCount
}

func (ctx *StratumContext) GetErrorCount() int64 {
	return ctx.errorCount
}

func (ctx *StratumContext) GetContext() context.Context {
	return ctx.ctx
}

func (ctx *StratumContext) NotifyLock() {
	ctx.notifyMutex.Lock()
}

func (ctx *StratumContext) NotifyUnlock() {
	ctx.notifyMutex.Unlock()
}

func (ctx *StratumContext) SetStartNotifyTs(startNotifyTs int64) *StratumContext {
	ctx.startNotifyTs = startNotifyTs
	return ctx
}

func (ctx *StratumContext) SetLatestNotifyTs(latestNotifyTs int64) *StratumContext {
	ctx.latestNotifyTs = latestNotifyTs
	return ctx
}

func (ctx *StratumContext) SetLatestDifficulty(latestDifficulty uint64) *StratumContext {
	ctx.latestDifficulty = latestDifficulty
	return ctx
}

func (ctx *StratumContext) SetLatestNotifyJobHeight(latestHeight int32) *StratumContext {
	ctx.latestNotifyJobHeight = latestHeight
	return ctx
}

func (ctx *StratumContext) SetLatestNotifyJobIndex(latestIndex int) *StratumContext {
	ctx.latestNotifyJobIndex = latestIndex
	return ctx
}

func (ctx *StratumContext) SetAuthorizeWorkerName(workerName string) *StratumContext {
	ctx.workerName = workerName
	return ctx
}

func (ctx *StratumContext) SetAuthorizeExtName(extName string) *StratumContext {
	ctx.extName = extName
	return ctx
}

func (ctx *StratumContext) SetAuthorizeUserName(userName string) *StratumContext {
	ctx.userName = userName
	return ctx
}

func (ctx *StratumContext) SetAuthorizePassword(password string) *StratumContext {
	ctx.password = password
	return ctx
}

func (ctx *StratumContext) SetAuthorizeTs(ts int64) *StratumContext {
	ctx.authorizeTs = ts
	return ctx
}

func (ctx *StratumContext) SetSubscribeSessionID(sessionID string) *StratumContext {
	ctx.sessionID = sessionID
	return ctx
}

func (ctx *StratumContext) SetSubscribeUserAgent(userAgent string) *StratumContext {
	ctx.userAgent = userAgent
	return ctx
}

func (ctx *StratumContext) SetSubscribeTs(ts int64) *StratumContext {
	ctx.subscribeTs = ts
	return ctx
}

func (ctx *StratumContext) SetConnRemoteIP(ip string) *StratumContext {
	ctx.remoteIP = ip
	return ctx
}

func (ctx *StratumContext) SetConnLocalIP(ip string) *StratumContext {
	ctx.localIP = ip
	return ctx
}

func (ctx *StratumContext) SetConnAcceptedTs(ts int64) *StratumContext {
	ctx.acceptedTs = ts
	return ctx
}

func (ctx *StratumContext) SetConnClosedTs(ts int64) *StratumContext {
	ctx.closedTs = ts
	return ctx
}

func (ctx *StratumContext) SetSlideWindow(window *average.SlidingWindow) *StratumContext {
	ctx.slideWindow = window
	return ctx
}

func (ctx *StratumContext) SetLatestSubmitTs(ts int64) *StratumContext {
	ctx.latestSubmitTs = ts
	return ctx
}

func (ctx *StratumContext) IncrNotifyCount() *StratumContext {
	ctx.notifyCount++
	return ctx
}

func (ctx *StratumContext) IncrSubmitCount() *StratumContext {
	ctx.submitCount++
	return ctx
}

func (ctx *StratumContext) IncrAcceptCount() *StratumContext {
	ctx.acceptCount++
	return ctx
}

func (ctx *StratumContext) IncrRejectCount() *StratumContext {
	ctx.rejectCount++
	return ctx
}

func (ctx *StratumContext) IncrErrorCount() *StratumContext {
	ctx.errorCount++
	return ctx
}

func (ctx *StratumContext) Close() {
	ctx.cancelFunc()
}
