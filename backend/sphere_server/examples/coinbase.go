package main

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"math/big"
	"strings"
	"time"
)

func HToBe(v uint32) uint32 {
	return ((v & 0xff000000) >> 24) |
		((v & 0x00ff0000) >> 8) |
		((v & 0x0000ff00) << 8) |
		((v & 0x000000ff) << 24);
}

func main() {
	merkleRoot1()
	//extraNonce2()
	//merkleBranch()
	//merkleRoot()
	//prevHash2()
	//fmt.Println(len("00000000c248468108d8cbdb546d542beaffd8c6470c25aeb2b6d09fa7c24962"))
	//data := "00000020723780e9701f632f150cb0eef3e0525db254764022a1bde547000000000000004433ee5a08bac08b7dd9cb38c488e3f88c3acc3617ef6ab335a2a0c9b8044b494486625c3568301787399505"
	//fmt.Println(len(data))
	//buf, _ := hex.DecodeString(data)
	//fmt.Println(chainhash.DoubleHashH(buf).String())

	//data = "20000000e98037722f631f70eeb00c155d52e0f3407654b2e5bda1220000004700000000f4642d1803b3b40579049cb115173e34b19f6b7c6002a9a34c9d9b2bcb0a04b55c6286441730683505953987"
	//fmt.Println(len(data))
	//buf, _ = hex.DecodeString(data)
	//fmt.Println(chainhash.DoubleHashH(buf).String())

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


func extraNonce2() {
	extraNonce2 := "01000000"
	nonceBytes, _ := hex.DecodeString(extraNonce2)
	extraNonce2Num := binary.LittleEndian.Uint32(nonceBytes)
	fmt.Println(extraNonce2Num)
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, extraNonce2Num)
	fmt.Println(hex.EncodeToString(a), extraNonce2Num)

}

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
	//extraNonce1 := "58685a615a546463"
	extraNonce1 := "7349424b534c6866"
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

func littleEndian2(text string) string {
	slice := make([]string, 0)
	for i := 0; i < len(text); i = i + 2 {
		slice = append(slice, reverseString(string(text[i:i+2])))
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
	//hash:= "b582dd0f79d653ee665ed7cb1eabcfb20ba7b69ecf16f436005c8a141e008fb2"
	fmt.Println(littleEndian(reverseString(hash)))

}

func prevHash2() {
	// 89c2f63dfb970e5638aa66ae3b7404a8a9914ad80328e9fe0000000000000000
	// 00000000000000000328e9fea9914ad83b7404a838aa66aefb970e5689c2f63d

	h := "534c68667349424b"
	h = littleEndian(reverseString(h))
	fmt.Println(h)
	//
	//h = "0000000000000000002b634cb1c59764a2d5f3c1d7656441e49563bc40b7d97a"
	//hash, err := chainhash.NewHashFromStr(h)
	//fmt.Println(hash, err)
}


func merkleRoot1() {
	//fmt.Println(littleEndian(reverseString("b5cf146c03660c572979b58210fe1d68bade9ae2001831180000000000000000")))
	//fmt.Println(littleEndian(reverseString("bab0cd4528324df44f10655a1359a326a560effc3e3d21925d4ee779b992aa90")))

	//prevHash, _ := chainhash.NewHashFromStr("000000000000000000183118bade9ae210fe1d682979b58203660c57b5cf146c")
	//merkleRoot, _ := chainhash.NewHashFromStr("b992aa905d4ee7793e3d2192a560effc1359a3264f10655a28324df4bab0cd45")
	//time, _ := decodeTimestamp("5d138334")

	//merkleBranches := []string{
	//	"39ee6edf0d227c264e86ea75561702ce00a3cd4a8d1a4183b9b905cce603c187",
	//	"c806c1704ac817ebfca6f22fe43190800f9f92bce4c96ca518bd4a05677becb8",
	//	"9b577e03d39f59845d615fdceb35a0bfb4aab338f63bfa29c7049038d9a6d434",
	//	"1a9fc7088f2c81bf83709692a1c14c1fa3bba40bdb8aa5bdcf61ddd50e09c826",
	//	"42611ea96ac2fa76c1c45606468907dd581ce205eb54c3e7bdaab69d79020f7b",
	//	"3d5c5bf72f971bebd4c66945fd91d4218e9d268b3576ab4e610f7c1e12c7ce20",
	//	"82958e5b98640512b5d97d244b08cdce7d2e72fc05f4915065400cb4e50400c5",
	//	"d031b5187a374ea8f494025d9e686bae7e57abb6852ee155dec9c9d3342bed45",
	//	"eafbc05346f0974a2d0a139b933d85b71173f137d0a34ead8ffd9ce2f70a893c",
	//	"88061c3a3dc2e723d423f0916edc6308f52ca428d1e7230450d3105560efa0c6",
	//	"9396efb267791b12d90d8ad0ad61aa3269cf91e2a240f7ad48f00daceba0645b",
	//	"faa6f07271ca8b2b542c91ecb8071ec8265e5bcde5e211037f8dd329bbcae1b8",
	//}

	//coinBaseBuf := new(bytes.Buffer)
	//coinBaseBuf.Write([]byte("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff2d03f1e2080876ec7df4057aab150a2f6c616b65706f6f6c2f086637653039663732"))
	//coinBaseBuf.Write([]byte("4266644f74637678"))
	//coinBaseBuf.Write([]byte("01000000"))
	//coinBaseBuf.Write([]byte("ffffffff02f82f2a50000000001976a914c5f7273961fd259947f73901d67ea3d20e67543988ac0000000000000000066a24aa21a9ed00000000"))

	//txBytes, err := hex.DecodeString("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff2d0385e308086ba51c853fc6ab150a2f6c616b65706f6f6c2f0866376530396637327a43457a7a65537900000000ffffffff02294f7551000000001976a914c5f7273961fd259947f73901d67ea3d20e67543988ac0000000000000000066a24aa21a9ed00000000")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//coinBaseMsgTx := wire.NewMsgTx(wire.TxVersion)
	////coinBaseMsgTx.BtcDecode(bytes.NewReader(txBytes), wire.BIP0111Version, wire.BaseEncoding)
	//err = coinBaseMsgTx.Deserialize(bytes.NewReader(txBytes))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	////merkleRootStr := BuildMerkleRoot(coinBaseMsgTx.TxHash().String(), merkleBranches)
	//
	//fmt.Println(btcutil.NewTx(coinBaseMsgTx).Hash())

	////////

	//hash 计算


	//没问题
	prevHash, _ := chainhash.NewHashFromStr(littleEndian(reverseString("e3159942504424b2f8d365a97f91a64cc984b80200030f530000000000000000")))
	merkleRoot, _ := chainhash.NewHashFromStr(littleEndian(reverseString("860ff844268df268c5c4cef7eb76ce27a4719291f88a37653e5bdcaba18ce0e6")))
	time, _ := decodeTimestamp("5d13dfb5")
	header := wire.BlockHeader{
		Version:    0X20000000,
		PrevBlock:  *prevHash,
		MerkleRoot: *merkleRoot,
		Timestamp:  time,
		Bits:       0x1725fd03,
		Nonce:      0xec607e0b,
	}

	fmt.Println(*prevHash, *merkleRoot, header.BlockHash(),  time)
	h, _ := chainhash.NewHashFromStr(header.BlockHash().String())
	//merkle root 没问题
	hh, _ := chainhash.NewHashFromStr("00000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	//fmt.Println(littleEndian2(reverseString("a18ce0e63e5bdcabf88a3765a4719291eb76ce27c5c4cef7268df268860ff844")))
	fmt.Println(blockchain.HashToBig(h), blockchain.HashToBig(hh), blockchain.HashToBig(h).Cmp(new(big.Int).Div(blockchain.HashToBig(hh), new(big.Int).SetUint64(1))))


	//tx hash
	//txBytes, err := hex.DecodeString("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff2d03bde30808b7c14ce3d6dbab150a2f6c616b65706f6f6c2f0866376530396637326766517472474c7900000000ffffffff023b736351000000001976a914c5f7273961fd259947f73901d67ea3d20e67543988ac0000000000000000066a24aa21a9ed00000000")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//coinBaseMsgTx := wire.NewMsgTx(wire.TxVersion)
	//coinBaseMsgTx.BtcDecode(bytes.NewReader(txBytes), wire.BIP0111Version, wire.BaseEncoding)
	//err = coinBaseMsgTx.Deserialize(bytes.NewReader(txBytes))
	//if err != nil {
	//	fmt.Println(err)
	//}

	//txhash 没问题
	//fmt.Println("89008719af522083f99ed6e1084c8217651bb9c68c94b93aaa7643da06db5834")
	//fmt.Println(littleEndian2(reverseString("89008719af522083f99ed6e1084c8217651bb9c68c94b93aaa7643da06db5834")))
	//fmt.Println(littleEndian2(reverseString(coinBaseMsgTx.TxHash().String())))
	//fmt.Println(coinBaseMsgTx.TxHash().String())


	//
	////merkle root
	//merkleBranches := []string{
	//	"6bd88c2f1f2bcf239f9c2f19f9fb3c99fe0e79c504d34623c06ea710e638f611",
	//	"be99d3b6bcde49add70fe4585d5925175a6d7df10b15ee4f889ec7dcd814d244",
	//	"72b9a1f04030cc2d47afb8a09a2e9c816c436deebd8faf6bac85da746f1f2170",
	//	"b1c00684ad99434bb73058451fb95dced403c8f96cfe4ea2687a56a2371a9a76",
	//	"2e8c9097bd53234d1ed7931ff324568f9b0f29a1e13f9163a3f4e5682e1f9c2e",
	//	"1475d8d48efe8d96c558384e044cb6d27b352333e87e4904dc526eaeea923da6",
	//	"66ef7459fa4557b4f7bce758aabac629636620bb331624eb98c76f1896dc36ce",
	//	"8bb9c0766089cb51b619cd1452f5c42809fb6df2ed84ecc3cf4ea09a90d77ba3",
	//	"db7da3717ff61cebee30cd0d3a3e8862e2689a046e46bf3af33ce29b86e7618c",
	//	"991c4bf23764d2804c08f702a3add350c74bae9349f7e21876f66f357955f1f9",
	//	"fda72457ce6e1a928032512067f380ded2b7dc712688002a70b114e20fc1ff0a",
	//	"a1b89768df8887ac928347bbfaa4ccd290434740a01fd45fc0a909e73a4cb787",
	//}
	//
	//merkleRootStr := BuildMerkleRoot(coinBaseMsgTx.TxHash().String(), merkleBranches)
	//fmt.Println(merkleRootStr)

	//merkleRootStr1 := BuildMerkleRoot(littleEndian2(reverseString("7a357cfd49b7d36201933659dc597020e85c92adb68db51c6425e247b1939868")), []string{littleEndian2(reverseString("bd7a0b57475bc626205dcb65083504b0a6fed53708adffebf36f3488d95aeb96"))})
	//fmt.Println(littleEndian2(reverseString(merkleRootStr1)))



}





func merkleRoot() {

	//coinbase := "dff3754157e6555226403fcf96f82ae24c6b2cdbfd03ece42e1e365742d50357"
	//
	merkleBranches := []string{
		"8784ecb734dbb91b503478defe1b2807a4c7b978bd1d8e093d0bd072441815dd","705e5f619ca9c0cb8f599f3aafb4c13784d308d1d57949abc910e42a193b9266",
		"f3bb15e2598b370e5d8c9135be25371850e44d8497e695256f8131f5cc3d98cb","87d3793a04b700d5411bf31ef2ac9144b54407e1bd274bf8afd2196e9dfa9187",
		"699b581ed7e2e14314ea183d4f5388db785c4f987d61649deb176dae87f52955","789476293275766b5ce876640d909f65ce67830474c472538e2a742ab9b30fd7",
		"f877c7b4786ba6b3c0f0bdf3c789bd4e0cc2abf0ee09de32c54bbf8eb91c83ae","1f72ca3a8777f69488e3bfb86c79a8685b35f51c9ef02514f5be124b1164a3f0",
		"e464486c4cd48bf2cadddbd745e07d87f96fe6db9a7fb1e688fe8972cd3142e6","3424005f7a9a574c4199db05fd2581d63f366c8f95943c7f610de912a6938e8f",
		"2c13f6748b0f5cc8b28340d0ceca3a917acf8c40a61719df0eaa596a0f551a13","0a2fd11f2a1d08343ec7435d15669f51178b5e7c42c15f3bbc49a4911501d6d1",
	}

	coinBaseBuf := new(bytes.Buffer)
	coinBaseBuf.Write([]byte("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff2d03c2e208085a7a2008915eab150a2f6c616b65706f6f6c2f086637653039663732"))
	coinBaseBuf.Write([]byte("534c68667349424b"))
	coinBaseBuf.Write([]byte("00000000"))
	coinBaseBuf.Write([]byte("ffffffff02a63a164d000000001976a914c5f7273961fd259947f73901d67ea3d20e67543988ac0000000000000000066a24aa21a9ed00000000"))

	txBytes, err := hex.DecodeString(coinBaseBuf.String())
	if err != nil {
		fmt.Println(err)
	}

	coinBaseMsgTx := wire.NewMsgTx(wire.TxVersion)
	err = coinBaseMsgTx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(coinBaseMsgTx.TxHash(), coinBaseMsgTx.HasWitness(), coinBaseMsgTx.WitnessHash())
	merkleRootStr := BuildMerkleRoot(coinBaseMsgTx.TxHash().String(), merkleBranches)
	merkleRootStr1 := "1e008fb2005c8a14cf16f4360ba7b69e1eabcfb2665ed7cb79d653eeb582dd0f"
	fmt.Println(merkleRootStr, merkleRootStr1)
	merkleRoot, _ := chainhash.NewHashFromStr(merkleRootStr1)
	time, _ := decodeTimestamp("5d1252db")
	prevHash, _ := chainhash.NewHashFromStr("00000000000000000004456ebbcc46f86cda3c72fee0208a2cb9f1e5d4411460")
	//////fmt.Println(merkleRoot)
	//////
	header := wire.BlockHeader{
		Version:    0X20000000,
		PrevBlock:  *prevHash,
		MerkleRoot: *merkleRoot,
		Timestamp:  time,
		Bits:       0x1725fd03,
		Nonce:      0x47bf6308,
	}
	//////
	fmt.Println(*prevHash, *merkleRoot, header.BlockHash(),  time)
	//buffer := new(bytes.Buffer)
	//buffer1 := new(bytes.Buffer)
	//fmt.Println(header.Serialize(buffer), hex.EncodeToString(buffer.Bytes()))
	//fmt.Println(header.BtcEncode(buffer1, wire.BIP0037Version, wire.WitnessEncoding), hex.EncodeToString(buffer1.Bytes()))
	//head :=  &wire.BlockHeader{}
	//
	//txBytes, _ = hex.DecodeString("00000020601441d4e5f1b92c8a20e0fe723cda6cf846ccbb6e4504000000000000000000b28f001e148a5c0036f416cf9eb6a70bb2cfab1ecbd75e66ee53d6790fdd82b524e0115d03fd251712422c73")
	//err = head.Deserialize(bytes.NewReader(txBytes))
	//fmt.Println(err)
	//fmt.Println(head.Version, head.PrevBlock.String(), head.Timestamp.String(), head.Bits, head.Nonce, head.BlockHash())

	//coinBaseTx := btcutil.NewTx(coinBaseMsgTx)
	//blockTxes := make([]*btcutil.Tx, 3)
	//blockTxes[0] = coinBaseTx
	//blockTxes[1] = coinBaseTx
	//blockTxes[2] = coinBaseTx
	//txHashes := make([]*chainhash.Hash, 3)
	//for i, transaction := range blockTxes {
	//	txHashes[i] = transaction.Hash()
	//}
	//merkles := blockchain.BuildMerkleTreeStore(blockTxes,false )
	//fmt.Println(merkles, merkles[len(merkles) -1], txHashes)
	//fmt.Println(merkles[len(merkles) - 1])
	//merkles = blockchain.BuildMerkleTreeStore(blockTxes, true)
	//fmt.Println(merkles[len(merkles) - 1])

}

func merkleBranch() {
	var txHashes1 = []string{
		"4d91b2aaf0b49c489c8d9a4bb545c00bdde5d8cd58a3cb06d09fa706beb2dc36", "fea89d9f9ea1946f30b969d0d662ca8b1ddd64ef05b60b11a4349098648f0ece",
		"4d91b2aaf0b49c489c8d9a4bb545c00bdde5d8cd58a3cb06d09fa706beb2dc36", "fea89d9f9ea1946f30b969d0d662ca8b1ddd64ef05b60b11a4349098648f0ece",
	}
	fmt.Println(makeMerkleBranch(txHashes1), len(txHashes1))
	fmt.Println(makeMerkleBranch(txHashes1), len(txHashes1))
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
	var hashes = make([]string, len(txHashes))
	copy(hashes, txHashes)
	for len(hashes) > 1 {
		branchesList.PushBack(hashes[0])
		if len(hashes)%2 == 0 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}
		for i := 0; i < (len(hashes)-1)/2; i++ {
			hashes[i] = HashMerkleBranches(hashes[i*2+1], hashes[i*2+2])
		}
		hashes = hashes[:(len(hashes)-1)/2]
	}
	branchesList.PushBack(hashes[0])
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

func BuildMerkleRoot(coinbase string, merkleBranches []string) string {
	var tmp *chainhash.Hash
	for _, branch := range merkleBranches {
		branch, _ := chainhash.NewHashFromStr(littleEndian2(reverseString(branch)))
		//branch, _ := chainhash.NewHashFromStr(branch)
		if tmp == nil  {
			tmp, _ = chainhash.NewHashFromStr(coinbase)
		}
		//fmt.Println("----", tmp, branch)
		tmp = blockchain.HashMerkleBranches(tmp, branch)
		//fmt.Println(tmp)
	}
	//fmt.Println("=====", tmp)
	return tmp.String()
}

func decodeTimestamp(ts string) (time.Time, error) {
	timeBytes, err := hex.DecodeString(ts)
	if err != nil {
		return time.Unix(0, 0), err
	}
	timestamp := binary.BigEndian.Uint32(timeBytes)

	return time.Unix(int64(timestamp), 0), nil
}
