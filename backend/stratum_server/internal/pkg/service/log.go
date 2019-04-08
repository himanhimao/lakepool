package service

import "time"

type BlockLog struct {
	workerName string
	serverIP   string
	clientIP   string
	userName   string
	extName    string
	userAgent  string
	hostName   string
	pid        int32
	height     int32
	hash       string
}

type LogService interface {
	Init() error
	AddBlockLog(coinType string, t time.Time, blockLog *BlockLog) (bool, error)
}

func NewBlockLog() *BlockLog {
	return &BlockLog{}
}

func (l *BlockLog) GetWorkerName() string {
	return l.workerName
}

func (l *BlockLog) GetServerIP() string {
	return l.serverIP
}

func (l *BlockLog) GetClientIP() string {
	return l.clientIP
}

func (l *BlockLog) GetUserName() string {
	return l.userName
}

func (l *BlockLog) GetExtName() string {
	return l.extName
}

func (l *BlockLog) GetUserAgent() string {
	return l.userAgent
}

func (l *BlockLog) GetHostName() string {
	return l.hostName
}

func (l *BlockLog) GetPid() int32 {
	return l.pid
}

func (l *BlockLog) GetHeight() int32 {
	return l.height
}

func (l *BlockLog) GetHash() string {
	return l.hash
}

func (l *BlockLog) SetWorkerName(workerName string) *BlockLog {
	l.workerName = workerName
	return l
}

func (l *BlockLog) SetServerIP(ip string) *BlockLog {
	l.serverIP = ip
	return l
}

func (l *BlockLog) SetClientIP(ip string) *BlockLog {
	l.clientIP = ip
	return l
}

func (l *BlockLog) SetUserName(userName string) *BlockLog {
	l.userName = userName
	return l
}

func (l *BlockLog) SetExtName(extName string) *BlockLog {
	l.extName = extName
	return l
}

func (l *BlockLog) SetUserAgent(userAgent string) *BlockLog {
	l.userAgent = userAgent
	return l
}

func (l *BlockLog) SetHostName(hostName string) *BlockLog {
	l.hostName = hostName
	return l
}

func (l *BlockLog) SetHeight(height int32) *BlockLog {
	l.height = height
	return l
}

func (l *BlockLog) SetPid(pid int32) *BlockLog {
	l.pid = pid
	return l
}

func (l *BlockLog) SetHash(hash string) *BlockLog {
	l.hash = hash
	return l
}
