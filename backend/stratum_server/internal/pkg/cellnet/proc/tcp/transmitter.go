package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/util"
	"errors"
	"io"
	"net"
)

var (
	ErrSocketFlood = errors.New("Socket flood detected from")
)

type TCPMessageTransmitter struct {
}

type socketOpt interface {
	ApplySocketReadTimeout(conn net.Conn, callback func())
	ApplySocketWriteTimeout(conn net.Conn, callback func())
}

func (TCPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {
	reader, ok := ses.Raw().(io.Reader)
	if !ok || reader == nil {
		return nil, nil
	}
	opt := ses.Peer().(socketOpt)

	if conn, ok := ses.Raw().(net.Conn); ok {
		// 有读超时时，设置超时
		opt.ApplySocketReadTimeout(conn, func() {
			msg, err = util.RecvPacket(reader)
			return
		})
	}
	return
}

func (TCPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) (err error) {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := ses.Peer().(socketOpt)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(ses.Raw().(net.Conn), func() {
		err = util.SendPacket(writer, ses.(cellnet.ContextSet), msg)
	})

	return
}
