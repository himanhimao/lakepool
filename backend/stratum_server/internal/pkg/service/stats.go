package service

import "time"

type ShareLog struct {
	workerName    string
	serverIP      string
	clientIP      string
	userName      string
	extName       string
	userAgent     string
	hostName      string
	pid           int32
	height        int32
	compute_power float64
	isRight       bool
}

type StatsService interface {
	Init() error
	AddShareLog(coinType string, t time.Time, shareLog *ShareLog) (bool, error)
}

func NewShareLog() *ShareLog {
	return &ShareLog{}
}

func (l *ShareLog) GetWorkerName() string {
	return l.workerName
}

func (l *ShareLog) GetServerIP() string {
	return l.serverIP
}

func (l *ShareLog) GetClientIP() string {
	return l.clientIP
}

func (l *ShareLog) GetUserName() string {
	return l.userName
}

func (l *ShareLog) GetExtName() string {
	return l.extName
}

func (l *ShareLog) GetUserAgent() string {
	return l.userAgent
}

func (l *ShareLog) GetHostName() string {
	return l.hostName
}

func (l *ShareLog) GetPid() int32 {
	return l.pid
}

func (l *ShareLog) GetHeight() int32 {
	return l.height
}

func (l *ShareLog) GetComputePower() float64 {
	return l.compute_power
}

func (l *ShareLog) GetIsRight() bool {
	return l.isRight
}

func (l *ShareLog) SetWorkerName(workerName string) *ShareLog {
	l.workerName = workerName
	return l
}

func (l *ShareLog) SetServerIP(ip string) *ShareLog {
	l.serverIP = ip
	return l
}

func (l *ShareLog) SetClientIP(ip string) *ShareLog {
	l.clientIP = ip
	return l
}

func (l *ShareLog) SetUserName(userName string) *ShareLog {
	l.userName = userName
	return l
}

func (l *ShareLog) SetExtName(extName string) *ShareLog {
	l.extName = extName
	return l
}

func (l *ShareLog) SetUserAgent(userAgent string) *ShareLog {
	l.userAgent = userAgent
	return l
}

func (l *ShareLog) SetHostName(hostName string) *ShareLog {
	l.hostName = hostName
	return l
}

func (l *ShareLog) SetHeight(height int32) *ShareLog {
	l.height = height
	return l
}

func (l *ShareLog) SetPid(pid int32) *ShareLog {
	l.pid = pid
	return l
}

func (l *ShareLog) SetComputePower(computePower float64) *ShareLog {
	l.compute_power = computePower
	return l
}

func (l *ShareLog) SetIsRight(isRight bool) *ShareLog {
	l.isRight = isRight
	return l
}