package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/service"
	log "github.com/sirupsen/logrus"
	"math/big"
)

const (
	PlaceHolder           = 0xEE
	BaseDifficulty uint64 = 4294967296
	MaxHash string = "00000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
)

type Coin struct {
	rpcClient *RpcClient
}

func NewCoin() *Coin {
	return &Coin{}
}

func NewBTCoinWithArgs(client *RpcClient) *Coin {
	return &Coin{rpcClient: client}
}

func (c *Coin) SetRPCClient(client *RpcClient) *Coin {
	c.rpcClient = client
	return c
}

func (c *Coin) IsValidAddress(address string, isUsedTestNet bool) bool {
	var err error
	var result bool = true
	if isUsedTestNet {
		_, err = btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)

	} else {
		_, err = btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	}
	if err != nil {
		result = false
	}

	return result
}

func (c *Coin) GetLatestStratumJob(ctx *service.Register) (*service.StratumJobPart, []*service.Transaction, error) {
	var getBlockTemplateParams [1]interface{}
	getBlockTemplateParams[0] = map[string][]string{
		"rules": {"segwit"},
	}
	gbtBlockTemplate, err := c.rpcClient.GetBlockTemplate(getBlockTemplateParams)
	if err != nil {
		return nil, nil, err
	}

	payoutAddressStr := ctx.PayoutAddress
	version := gbtBlockTemplate.Version
	gbtPrevBlockHash := gbtBlockTemplate.PreviousBlockHash
	coinBaseValue := gbtBlockTemplate.CoinBaseValue
	nBits := gbtBlockTemplate.Bits
	witnessCommitment := gbtBlockTemplate.DefaultWitnessCommitment
	height := gbtBlockTemplate.Height
	blockTemplateTransactions := gbtBlockTemplate.Transactions
	minTimeTs := gbtBlockTemplate.MinTime
	curTimeTs := gbtBlockTemplate.CurTime
	prevBlockHash := littleEndianUint32(reverseString(gbtPrevBlockHash))
	placeHolderSize := ctx.ExtraNonce1Length + ctx.ExtraNonce2Length
	placeHolders := genPlaceHolders(placeHolderSize)

	payoutAddress, err := btcutil.DecodeAddress(payoutAddressStr, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, err
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	pkScript, err := txscript.PayToAddrScript(payoutAddress)
	coinBaseSignatureScript, err := generateCoinBaseScript(height, ctx.Id, ctx.PoolTag)

	if err != nil {
		return nil, nil, err
	}

	//var witnessNonce [blockchain.CoinbaseWitnessDataLen]byte
	tx.AddTxIn(&wire.TxIn{
		// Coinbase transactions have no inputs, so previous outpoint is
		// zero hash and max index.
		PreviousOutPoint: *wire.NewOutPoint(&chainhash.Hash{},
			wire.MaxPrevOutIndex),
		SignatureScript: addPlaceHolders(coinBaseSignatureScript, placeHolders),
		Sequence:        wire.MaxTxInSequenceNum,
		//Witness: wire.TxWitness{witnessNonce[:]},
	})

	tx.AddTxOut(&wire.TxOut{
		Value:    coinBaseValue,
		PkScript: pkScript,
	})

	if len(witnessCommitment) > 0 {
		witnessCommitmentDecodebuffer := new(bytes.Buffer)
		hex.Decode([]byte(witnessCommitment), witnessCommitmentDecodebuffer.Bytes())
		witnessScript := append(blockchain.WitnessMagicBytes, witnessCommitmentDecodebuffer.Bytes()...)
		// Finally, create the OP_RETURN carrying witness commitment
		// output as an additional output within the coinbase.
		commitmentOutput := &wire.TxOut{
			Value:    0,
			PkScript: witnessScript,
		}
		tx.AddTxOut(commitmentOutput)
	}

	coinBaseTx := btcutil.NewTx(tx)
	buf := new(bytes.Buffer)
	coinBaseTx.MsgTx().Serialize(buf)

	coinBaseTxHex := hex.EncodeToString(buf.Bytes())
	coinBase1, coinBase2, err := splitCoinBaseHex(coinBaseTxHex, placeHolders)


	if err != nil {
		return nil, nil, err
	}

	var merkleBranch []string
	var txHashes []*chainhash.Hash
	var transactions []*service.Transaction

	if len(blockTemplateTransactions) > 0 {
		txHashes = make([]*chainhash.Hash, len(blockTemplateTransactions))
		transactions = make([]*service.Transaction, len(blockTemplateTransactions))

		for i, txTemplate := range blockTemplateTransactions {
			tx := wire.NewMsgTx(wire.TxVersion)
			txBytes, _ := hex.DecodeString(txTemplate.Data)
			err := tx.Deserialize(bytes.NewReader(txBytes))
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Warningln("make block deserialize error")
				continue
			}
			buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSizeStripped()))
			tx.SerializeNoWitness(buf)
			txHashes[i] = btcutil.NewTx(tx).Hash()
			txHexData := hex.EncodeToString(buf.Bytes())
			transactions[i] = service.NewBlockTransactionPart(tx.TxHash().String(), txHexData)
		}
		merkleBranch = makeMerkleBranch(txHashes)
	}

	stratumJobPart := service.NewStratumJobPart()
	stratumJobPart.CoinBase1 = coinBase1
	stratumJobPart.CoinBase2 = coinBase2
	stratumJobPart.MerkleBranch = merkleBranch
	stratumJobPart.NBits = nBits
	stratumJobPart.Version = fmt.Sprintf("%x", version)
	stratumJobPart.PrevHash = prevBlockHash

	stratumJobMetaPart := service.NewStratumJobMetaPart()
	stratumJobMetaPart.Height = height
	stratumJobMetaPart.CurTimeTs = curTimeTs
	stratumJobMetaPart.MinTimeTs = minTimeTs
	stratumJobPart.Meta = stratumJobMetaPart

	return stratumJobPart, transactions, nil
}

func (c *Coin) MakeBlock(ctx *service.Register, header *service.BlockHeaderPart, base *service.BlockCoinBasePart, transactions []*service.Transaction) (*service.Block, error) {
	var block wire.MsgBlock

	blockVersion, err := decodeVersion(header.Version)
	if err != nil {
		return nil, err
	}

	timestamp, err := decodeTimestamp(header.NTime)
	if err != nil {
		return nil, err
	}

	prevHash, err := decodeHash(header.PrevHash)
	if err != nil {
		return nil, err
	}

	nonce, err := decodeNonce(header.Nonce)
	if err != nil {
		return nil, err
	}

	bits, err := decodeBits(header.NBits)
	if err != nil {
		return nil, err
	}

	coinBaseBuf := new(bytes.Buffer)
	coinBaseBuf.Write([]byte(base.CoinBase1))
	coinBaseBuf.Write([]byte(base.ExtraNonce1))
	coinBaseBuf.Write([]byte(base.ExtraNonce2))
	coinBaseBuf.Write([]byte(base.CoinBase2))

	log.WithFields(log.Fields{
		"version":         header.Version,
		"prev_block_hash": header.PrevHash,
		"timestamp":       header.NTime,
		"bits":            header.NBits,
		"nonce":           header.Nonce,
		"coinbase_1":      base.CoinBase1,
		"coinbase_2":      base.CoinBase2,
		"extranonce_1":    base.ExtraNonce1,
		"extranonce_2":    base.ExtraNonce2,
	}).Debugln("make block source data")

	txBytes, err := hex.DecodeString(coinBaseBuf.String())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Errorln("decode coinbase buf  error")
		return nil, err
	}

	coinBaseMsgTx := wire.NewMsgTx(wire.TxVersion)
	err = coinBaseMsgTx.DeserializeNoWitness(bytes.NewReader(txBytes))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Errorln("decode coinbase tx error")
		return nil, err
	}

	var merkleBranch []string
	coinBaseTx := btcutil.NewTx(coinBaseMsgTx)
	blockTxes := make([]*btcutil.Tx, len(transactions)+1)
	blockTxes[0] = coinBaseTx

	txHashes := make([]*chainhash.Hash, len(transactions))
	for i, transaction := range transactions {
		tx, err := decodeTransaction(transaction.Data)
		if err != nil {
			log.Errorln("decode transaction error %s", err.Error())
			return nil, err
		}
		txHashes[i] = tx.Hash()
		blockTxes[i+1] = tx
	}
	merkleBranch = makeMerkleBranch(txHashes)
	merkleRootStr :=  buildMerkleRoot(coinBaseTx.Hash().String(), merkleBranch)
	merkleRoot, _:= chainhash.NewHashFromStr(merkleRootStr)
	block.Header = wire.BlockHeader{
		Version:    blockVersion,
		PrevBlock:  prevHash,
		MerkleRoot: *merkleRoot,
		Timestamp:  timestamp,
		Bits:       bits,
		Nonce:      nonce,
	}

	blockHeaderHash := block.Header.BlockHash()

	log.WithFields(log.Fields{
		"target_version":     blockVersion,
		"target_timestamp":   timestamp.Unix(),
		"target_prevHash":    prevHash.String(),
		"target_nonce":       nonce,
		"target_bits":        bits,
		"coinbase_tx_hash":   coinBaseTx.Hash().String(),
		"coinbase_tx_witness_hash": coinBaseMsgTx.WitnessHash().String(),
		"merkle_root": merkleRoot.String(),
		"block_hash":         blockHeaderHash.String(),
		"block_tx_0":         blockTxes[0].Hash().String(),
		"block_tx_len":       len(blockTxes),
		"merkle_branch":      merkleBranch,
		"coinbase_hex":       coinBaseBuf.String(),
	}).Debugln("make block target data")

	for _, tx := range blockTxes {
		if err := block.AddTransaction(tx.MsgTx()); err != nil {
			return nil, err
		}
	}
	utilBlock := btcutil.NewBlock(&block)
	blockData, _ := utilBlock.Bytes()
	b := service.NewBlock(blockHeaderHash.String(), hex.EncodeToString(blockData))
	return b, nil
}

func (c *Coin) SubmitBlock(data string) (bool, error) {
	var submitBlockParams [1]interface{}
	submitBlockParams[0] = data
	result, err := c.rpcClient.SubmitBlock(submitBlockParams)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (c *Coin) IsSolveHash(hash string, targetDifficulty *big.Int) (bool, error) {
	hashPtr, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return false, err
	}

	maxHash, _ := chainhash.NewHashFromStr(MaxHash)
	targetValue := new(big.Int).Div(blockchain.HashToBig(maxHash),  targetDifficulty)

	if blockchain.HashToBig(hashPtr).Cmp(targetValue) <= 0 {
		return true, nil
	}
	return false, nil
}

func (c *Coin) GetTargetDifficulty(bitsHex string) (*big.Int, error) {
	bits, err := decodeBits(bitsHex)
	if err != nil {
		return nil, err
	}
	return blockchain.CompactToBig(bits), nil
}

func (c *Coin) CalculateShareComputePower(targetDifficulty *big.Int) (*big.Int, error) {
	//TODO check variable
	baseDifficultyBig := new(big.Int).SetUint64(BaseDifficulty)
	return new(big.Int).Mul(baseDifficultyBig, targetDifficulty), nil
}

func (c *Coin) GetNewBlockHeight() (int, error) {
	miningInfo, err := c.rpcClient.GetMiningInfo()
	if err != nil {
		return 0, err
	}
	return miningInfo.Blocks + 1, nil
}
