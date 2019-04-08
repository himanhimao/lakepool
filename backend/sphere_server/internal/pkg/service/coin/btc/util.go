package btc

import (
	"strings"
	"encoding/hex"
	"bytes"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"container/list"
	"errors"
	"time"
	"encoding/binary"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/blockchain"
	"strconv"
)

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func littleEndian(text string) string {
	slice := make([]string, 0)
	for i := 0; i < len(text); i = i + 8 {
		slice = append(slice, reverseString(string(text[i:i+8])))
	}
	return strings.Join(slice, "")
}

func splitCoinBaseHex(coinBaseHex string, placeHolders []byte) (string, string, error) {
	coinBaseData := strings.Split(coinBaseHex, hex.EncodeToString(placeHolders))
	if len(coinBaseData) != 2 {
		return "", "", errors.New("invalid coinbase hex.")
	}
	return coinBaseData[0], coinBaseData[1], nil
}

func generateCoinBase(height int32, registerId string, poolTag string) []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte(strconv.Itoa(int(height))))
	timeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeBuf, uint64(time.Now().UnixNano()))
	buf.Write(timeBuf)
	buf.Write([]byte(registerId))
	buf.Write([]byte(poolTag))
	return buf.Bytes()
}

func genPlaceHolders(size int) []byte {
	buf := make([]byte, 1)
	buf[0] = byte(PlaceHolder)
	return bytes.Repeat(buf, size)
}

func addPlaceHolders(value []byte, placeHolders []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(value)
	buf.Write(placeHolders)
	return buf.Bytes()
}

func makeMerkleBranch(txHashes []*chainhash.Hash) []string {
	branchesList := list.New()
	for len(txHashes) > 1 {
		branchesList.PushBack(txHashes[0])
		if len(txHashes)%2 == 0 {
			txHashes = append(txHashes, txHashes[len(txHashes)-1])
		}
		for i := 0; i < (len(txHashes)-1)/2; i++ {
			txHashes[i] = blockchain.HashMerkleBranches(txHashes[i*2+1], txHashes[i*2+2])
		}
		txHashes = txHashes[:(len(txHashes)-1)/2]
	}
	branchesList.PushBack(txHashes[0])
	branches := make([]string, branchesList.Len())
	i := 0
	for e := branchesList.Front(); e != nil; e = e.Next() {
		branches[i] = e.Value.(*chainhash.Hash).String()
		i++
	}
	return branches
}

func decodeVersion(version string) (int32, error) {
	versionBytes, err := hex.DecodeString(version)
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint32(versionBytes)
	return int32(v), err
}

func decodeTimestamp(ts string) (time.Time, error) {
	timeBytes, err := hex.DecodeString(ts)
	if err != nil {
		return time.Unix(0, 0), err
	}
	timestamp := binary.BigEndian.Uint32(timeBytes)

	return time.Unix(int64(timestamp), 0), nil
}

func decodeHash(hash string) (chainhash.Hash, error) {
	var chainHash chainhash.Hash

	hashStr := littleEndian(reverseString(hash))
	hashPtr, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		return chainHash, err
	}
	return *hashPtr, nil
}

func decodeNonce(nonce string) (uint32, error) {
	nonceBytes, err := hex.DecodeString(nonce)
	if err != nil {
		return 0, err
	}
	n := binary.BigEndian.Uint32(nonceBytes)
	return n, nil
}

func decodeBits(bits string) (uint32, error) {
	bitsBytes, err := hex.DecodeString(bits)
	if err != nil {
		return 0, err
	}
	b := binary.BigEndian.Uint32(bitsBytes)
	return b, nil
}

func decodeTransaction(data string) (*btcutil.Tx, error) {
	tx, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	originTx := wire.NewMsgTx(wire.TxVersion)
	err = originTx.Deserialize(bytes.NewReader(tx))
	if err != nil {
		return nil, err
	}
	return btcutil.NewTx(originTx), nil
}
