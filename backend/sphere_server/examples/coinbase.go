package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"encoding/hex"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/blockchain"
	"container/list"
)

func HToBe(v uint32) uint32 {
	return ((v & 0xff000000) >> 24) |
		((v & 0x00ff0000) >> 8) |
		((v & 0x0000ff00) << 8) |
		((v & 0x000000ff) << 24);
}

func main() {

	fmt.Println(len("00000000c248468108d8cbdb546d542beaffd8c6470c25aeb2b6d09fa7c24962"))
	data := "00000020723780e9701f632f150cb0eef3e0525db254764022a1bde547000000000000004433ee5a08bac08b7dd9cb38c488e3f88c3acc3617ef6ab335a2a0c9b8044b494486625c3568301787399505"
	fmt.Println(len(data))
	buf, _ := hex.DecodeString(data)
	fmt.Println(chainhash.DoubleHashH(buf).String())

	data = "20000000e98037722f631f70eeb00c155d52e0f3407654b2e5bda1220000004700000000f4642d1803b3b40579049cb115173e34b19f6b7c6002a9a34c9d9b2bcb0a04b55c6286441730683505953987"
	fmt.Println(len(data))
	buf, _ = hex.DecodeString(data)
	fmt.Println(chainhash.DoubleHashH(buf).String())

	//var block wire.MsgBlock
	//err := block.Header.Deserialize(bytes.NewReader(buf))
	//fmt.Println("err", err)
	//fmt.Println("version", block.Header.Version, block.Header.PrevBlock.String())

	//Version()
	//Nonce()
	//prevHash2()
	//CoinBase6()
	//Coinbase2()
	//fmt.Println(hex.EncodeToString([]byte("336568448")))
	// //Coinbase2()
	//c := fmt.Sprintf("%x", 0x23)
	//fmt.Println(len(c), c)
	//fmt.Println(hex.EncodeToString([]byte("cswd")))
	//fmt.Printf("%X", "##")
	//fmt.Printf("%c", string(0x23))
	//a := makeMerkleBranch([]string{"c284853b65e7887c5fd9b635a932e2e0594d19849b22914a8e6fb180fea0954f", "c284853b65e7887c5fd9b635a932e2e0594d19849b22914a8e6fb180fea0954f", "28b1a5c2f0bb667aea38e760b6d55163abc9be9f1f830d9969edfab902d17a0f", "67878210e268d87b4e6587db8c6e367457cea04820f33f01d626adbe5619b3dd",  })
	//fmt.Println(a)

	//fmt.Println(HashMerkleBranches("28b1a5c2f0bb667aea38e760b6d55163abc9be9f1f830d9969edfab902d17a0f", "67878210e268d87b4e6587db8c6e367457cea04820f33f01d626adbe5619b3dd" ))
	//.Println(decodeBits1("172e6f88"))
	//
	//
	//tx, err  := hex.DecodeString("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff020000ffffffff0100e1f505000000002200201863143c14c5166804bd19203356da136c985678cd4d27a1b8c632960490326200000000")
	//fmt.Println(err)
	//fmt.Println(err)
	//fmt.Printf("#%v", originTx.TxOut[0].Value)
	//01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1366326332613538642f6c616b65706f6f6c2feeffffffff025e290913000000001976a914c5f7273961fd259947f73901d67ea3d20e67543988ac0000000000000000066a24aa21a9ed00000000
	//01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1330656365616262632f6c616b65706f6f6c2feeffffffff025e290913000000001976a914c5f7273961fd259947f73901d67ea3d20e67543988ac0000000000000000066a24aa21a9ed00000000
}

//func Coinbase3() {
//	addressStr := "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sl5k7"
//	address, err := btcutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	script, err := txscript.PayToAddrScript(address)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	//fmt.Printf("Script Hex: %x\n", script)
//	//poolCoinbaseInfo := "/BTC.COM/";
//	//tx, err := txscript.NewScriptBuilder().AddInt64(int64(uint64(0))).AddData([]byte(poolCoinbaseInfo)).
//	//	Script()
//	//fmt.Printf("Script Hex: %x-%s \n", tx, err)
//	originTx := wire.NewMsgTx(wire.TxVersion)
//	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
//	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
//

//}

func Nonce() {
	nonceHex := "2b22c815"
	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		fmt.Println(nonceBytes)
	}
	fmt.Println(binary.BigEndian.Uint32(nonceBytes))
}

func CoinBase4() {
	coinBase1 := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1e30333534393366352f6c616b65706f6f6c2f"
	coinBase2 := "ffffffff025e290913000000001976a914c5f7273961fd259947f73901d67ea3d20e67543988ac0000000000000000066a24aa21a9ed00000000"
	extraNonce1 := "58685a615a546463"
	extraNonce2 := "00000000"
	coinBaseBuf := new(bytes.Buffer)
	coinBaseBuf.Write([]byte(coinBase1))
	coinBaseBuf.Write([]byte(extraNonce1))
	coinBaseBuf.Write([]byte(extraNonce2))
	coinBaseBuf.Write([]byte(coinBase2))

	originTx := wire.NewMsgTx(wire.TxVersion)
	fmt.Println(coinBaseBuf.String())
	tx, err := hex.DecodeString(coinBaseBuf.String())

	fmt.Println(err)
	err = originTx.Deserialize(bytes.NewReader(tx))
	fmt.Println(err)
}

func CoinBase5() {
	transaction := "01000000013faf73481d6b96c2385b9a4300f8974b1b30c34be30000c7dcef11f68662de4501000000db00483045022100f9881f4c867b5545f6d7a730ae26f598107171d0f68b860bd973dbb855e073a002207b511ead1f8be8a55c542ce5d7e91acfb697c7fa2acd2f322b47f177875bffc901483045022100a37aa9998b9867633ab6484ad08b299de738a86ae997133d827717e7ed73d953022011e3f99d1bd1856f6a7dc0bf611de6d1b2efb60c14fc5931ba09da01558757f60147522102632178d046673c9729d828cfee388e121f497707f810c131e0d3fc0fe0bd66d62103a0951ec7d3a9da9de171617026442fcd30f34d66100fab539853b43f508787d452aeffffffff0240420f000000000017a9148d57003ecbaa310a365f8422602cc507a702197e87806868a90000000017a9148ce5408cfeaddb7ccb2545ded41ef478109454848700000000"
	tx, err := hex.DecodeString(transaction)
	originTx := wire.NewMsgTx(wire.TxVersion)
	err = originTx.Deserialize(bytes.NewReader(tx))
	fmt.Println(err)
}

func CoinBase6() {
	timeHex := "5c5efbd4"
	timeBytes, err := hex.DecodeString(timeHex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(binary.BigEndian.Uint32(timeBytes))

}

func createCoinbaseTx(params *chaincfg.Params, coinbaseScript []byte, nextBlockHeight int32, addr btcutil.Address) (*btcutil.Tx, error) {
	// Create the script to pay to the provided payment address if one was
	// specified.  Otherwise create a script that allows the coinbase to be
	// redeemable by anyone.
	var pkScript []byte
	if addr != nil {
		var err error
		pkScript, err = txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		scriptBuilder := txscript.NewScriptBuilder()
		pkScript, err = scriptBuilder.AddOp(txscript.OP_TRUE).Script()
		if err != nil {
			return nil, err
		}
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	tx.AddTxIn(&wire.TxIn{
		// Coinbase transactions have no inputs, so previous outpoint is
		// zero hash and max index.
		PreviousOutPoint: *wire.NewOutPoint(&chainhash.Hash{},
			wire.MaxPrevOutIndex),
		SignatureScript: coinbaseScript,
		Sequence:        wire.MaxTxInSequenceNum,
	})
	tx.AddTxOut(&wire.TxOut{
		Value:    blockchain.CalcBlockSubsidy(nextBlockHeight, params),
		PkScript: pkScript,
	})
	return btcutil.NewTx(tx), nil
}

func Coinbase2() {
	addressStr := "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sl5k7"
	address, err := btcutil.DecodeAddress(addressStr, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	script, err := txscript.PayToAddrScript(address)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Script Hex: %x\n", script)
	poolCoinbaseInfo := "/BTC.COM/";
	tx, err := txscript.NewScriptBuilder().AddInt64(int64(uint64(0))).AddData([]byte(poolCoinbaseInfo)).
		Script()
	fmt.Printf("Script Hex: %x-%s \n", tx, err)

	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	originTx.AddTxIn(txIn)
	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)

	buffer := new(bytes.Buffer)
	originTxHash := originTx.Serialize(buffer)

	fmt.Println(originTxHash)
	fmt.Println(fmt.Sprintf("%x", buffer))

	fmt.Println(fmt.Sprintf("%x", 0xEE))

}

func littleEndian(text string) string {
	slice := make([]string, 0)
	for i := 0; i < len(text); i = i + 8 {
		slice = append(slice, reverseString(string(text[i:i+8])))
	}
	return strings.Join(slice, "");
}

func Reverse(s []byte) []byte {
	L := len(s)
	t := make([]byte, L)
	for i, v := range s {
		t[L-1-i] = v
	}
	return t
}

func reverseString(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}

func test() {
	t := "4294967295"
	buf := new(bytes.Buffer)
	buf.WriteString(t)
	fmt.Println(hex.EncodeToString(buf.Bytes()))
}

/*
00000000000000000023ad41ca27a380047e3d14969e515582427a8eadb7814b
adb7814b82427a8e969e5155047e3d14ca27a3800023ad410000000000000000
*/
func coinbase1Main() {
	poolCoinbaseInfo := "/BTC.COM/";
	//blockVersion := 0;
	height := int32(500000)
	buf := new(bytes.Buffer)
	//registerId := "ccd8d64d0b4303df6890e0af69876e19"
	registerId := "ccd8d6"

	_ = binary.Write(buf, binary.LittleEndian, height)
	buf.Write([]byte(registerId))
	buf.Write([]byte(poolCoinbaseInfo))

	fmt.Println(len(buf.Bytes()), buf.Bytes())
}

func prevHash() {
	// previous block hash
	// we need to convert to little-endian
	// 00000000000000000328e9fea9914ad83b7404a838aa66aefb970e5689c2f63d
	// 89c2f63dfb970e5638aa66ae3b7404a8a9914ad80328e9fe0000000000000000

	/*
	for (int i = 0; i < 8; i++) {
	uint32_t a = *(uint32_t *)(BEGIN(prevHash_) + i * 4);
	a = HToBe(a);
	prevHashBeStr_ += HexStr(BEGIN(a), END(a));
	}
	*/
	hash := "00000000000000000328e9fea9914ad83b7404a838aa66aefb970e5689c2f63d"
	fmt.Println(littleEndian(reverseString(hash)))

}

func prevHash2() {
	// 89c2f63dfb970e5638aa66ae3b7404a8a9914ad80328e9fe0000000000000000
	// 00000000000000000328e9fea9914ad83b7404a838aa66aefb970e5689c2f63d

	h := "c81cf6f5c43e149d70a1ded7d0e8e24d94102ffe0014bec10000000000000000"
	h = littleEndian(reverseString(h))
	fmt.Println(h)
	//
	//h = "0000000000000000002b634cb1c59764a2d5f3c1d7656441e49563bc40b7d97a"
	//hash, err := chainhash.NewHashFromStr(h)
	//fmt.Println(hash, err)
}

func Version() {
	version := "20000000"
	//v, _ := strconv.Atoi(version)
	//
	versionBytes, err := hex.DecodeString(version)
	fmt.Println(err)
	fmt.Println(string(versionBytes))
	v := binary.BigEndian.Uint32(versionBytes)
	fmt.Println(v)
	s := fmt.Sprintf("%x", v)
	fmt.Println(s)
}

func makeMerkleBranch(txHashes []string) []string {
	branchesList := list.New()
	for len(txHashes) > 1 {
		branchesList.PushBack(txHashes[0])
		if len(txHashes)%2 == 0 {
			txHashes = append(txHashes, txHashes[len(txHashes)-1])
		}
		for i := 0; i < (len(txHashes)-1)/2; i++ {
			txHashes[i] = HashMerkleBranches(txHashes[i*2+1], txHashes[i*2+2])
		}
		txHashes = txHashes[:(len(txHashes)-1)/2]
	}
	branchesList.PushBack(txHashes[0])
	branches := make([]string, branchesList.Len())
	i := 0
	for e := branchesList.Front(); e != nil; e = e.Next() {
		branches[i] = e.Value.(string)
		i++
	}
	return branches
}

func HashMerkleBranches(left string, right string) string {
	// Concatenate the left and right nodes.
	var hash [chainhash.HashSize * 2]byte
	copy(hash[:chainhash.HashSize], left[:])
	copy(hash[chainhash.HashSize:], right[:])
	newHash := chainhash.DoubleHashH(hash[:])
	return newHash.String()
}
