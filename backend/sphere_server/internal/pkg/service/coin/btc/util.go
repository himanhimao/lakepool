package btc

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"strings"
	"time"
)

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func littleEndianUint32(text string) string {
	slice := make([]string, 0)
	for i := 0; i < len(text); i = i + 8 {
		slice = append(slice, reverseString(string(text[i:i+8])))
	}
	return strings.Join(slice, "")
}

func littleEndianUint8(text string) string {
	slice := make([]string, 0)
	for i := 0; i < len(text); i = i + 2 {
		slice = append(slice, reverseString(string(text[i:i+2])))
	}
	return strings.Join(slice, "");
}

func splitCoinBaseHex(coinBaseHex string, placeHolders []byte) (string, string, error) {
	coinBaseData := strings.Split(coinBaseHex, hex.EncodeToString(placeHolders))
	if len(coinBaseData) != 2 {
		return "", "", errors.New("invalid coinbase hex.")
	}
	return coinBaseData[0], coinBaseData[1], nil
}

func generateCoinBaseScript(height int32, registerId string, poolTag string) ([]byte, error) {
	return txscript.NewScriptBuilder().AddInt64(int64(height)).AddInt64(time.Now().UnixNano()).AddData([]byte(poolTag)).
		AddData([]byte(registerId)).Script()
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
	var hashes = make([]*chainhash.Hash, len(txHashes))
	copy(hashes, txHashes)

	branchesList := list.New()
	for len(hashes) > 1 {
		branchesList.PushBack(hashes[0])
		if len(hashes)%2 == 0 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}
		for i := 0; i < (len(hashes)-1)/2; i++ {
			hashes[i] = blockchain.HashMerkleBranches(hashes[i*2+1], hashes[i*2+2])
		}
		hashes = hashes[:(len(hashes)-1)/2]
	}
	branchesList.PushBack(hashes[0])
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

	hashStr := littleEndianUint32(reverseString(hash))
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

func incrExtraNonce2(extraNonce2 string, length int) string {
	buf := make([]byte, length)
	nonceBytes, _ := hex.DecodeString(extraNonce2)
	extraNonce2Num := binary.LittleEndian.Uint32(nonceBytes)
	binary.LittleEndian.PutUint32(buf, extraNonce2Num + 1)
	return hex.EncodeToString(buf)
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


func buildMerkleRoot(coinbase string, branches []string) string {
	var tmp *chainhash.Hash
	for _, branch := range branches {
		branch, _ := chainhash.NewHashFromStr(littleEndianUint8(reverseString(branch)))
		if tmp == nil  {
			tmp, _ = chainhash.NewHashFromStr(coinbase)
		}
		tmp = blockchain.HashMerkleBranches(tmp, branch)
	}
	return tmp.String()
}