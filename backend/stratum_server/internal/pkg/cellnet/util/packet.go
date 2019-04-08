package util

import (
	"bytes"
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/codec"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"
	"io"
)

const (
	RecvBufSize = 1
	MaxReadSize = 1024
)

var (
	ErrMsgFloodAttack = errors.New("flood attack")
)

func RecvPacket(reader io.Reader) (msg interface{}, err error) {
	msgBuf := new(bytes.Buffer)
	readSize := 0
	for {
		buf := make([]byte, RecvBufSize)

		if _, err := reader.Read(buf); err != nil {
			return nil, err
		}

		if buf[0] != '\n' {
			msgBuf.Write(buf)
		} else {
			break
		}
		readSize++

		if readSize > MaxReadSize {
			err = ErrMsgFloodAttack
			return nil, err
		}
	}
	//bufio.NewReader(reader)

	msg, _, err = codec.DecodeMessage(proto.MsgFullNameJsonRpcREQ, msgBuf.Bytes())
	if err != nil {
		return nil, err
	}
	return
}

func SendPacket(writer io.Writer, ctx cellnet.ContextSet, data interface{}) error {
	// 将用户数据转换为字节数组和消息ID
	msgData, meta, err := codec.EncodeMessage(data, ctx)

	if err != nil {
		return err
	}

	// 将数据写入Socket
	err = util.WriteFull(writer, msgData)
	if err != nil {
		return err
	}

	util.WriteFull(writer, []byte{'\n'})
	if err != nil {
		return err
	}

	// Codec中使用内存池时的释放位置
	if meta != nil {
		codec.FreeCodecResource(meta.Codec, msgData, ctx)
	}

	return nil
}
