package util

import (
	"encoding/hex"
	"encoding/binary"
)

func ConvertTs(hexTime string) (uint32, error) {
	var bytes []byte
	var err error
	var ts uint32
	if bytes, err = hex.DecodeString(hexTime); err != nil {
		return 0, err
	}
	ts = binary.BigEndian.Uint32(bytes)
	return ts, nil
}
