package btc

import (
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"bytes"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"encoding/hex"
	"github.com/btcsuite/btcd/blockchain"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/service"
	"fmt"
	"math/big"
)

const (
	PlaceHolder           = 0xEE
	BaseDifficulty uint64 = 4294967296
)

type BTCCoin struct {
	rpcClient *RpcClient
}

func NewBTCCoin() *BTCCoin {
	return &BTCCoin{}
}

func NewBTCoinWithArgs(client *RpcClient) *BTCCoin {
	return &BTCCoin{rpcClient: client}
}

func (c *BTCCoin) SetRPCClient(client *RpcClient) *BTCCoin {
	c.rpcClient = client
	return c
}

func (c *BTCCoin) IsValidAddress(address string, isUsedTestNet bool) bool {
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

func (c *BTCCoin) GetLatestStratumJob(registerId string, ctx *service.Register) (*service.StratumJobPart, []*service.BlockTransactionPart, error) {
	var getBlockTemplateParams [1]interface{}
	getBlockTemplateParams[0] = map[string][]string{
		"rules": []string{"segwit"},
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
	coinBaseSignatureScript := generateCoinBase(height, registerId, ctx.PoolTag)
	blockTemplateTransactions := gbtBlockTemplate.Transactions
	minTimeTs := gbtBlockTemplate.MinTime
	curTimeTs := gbtBlockTemplate.CurTime
	prevBlockHash := littleEndian(reverseString(gbtPrevBlockHash))
	placeHolderSize := ctx.ExtraNonce1Length + ctx.ExtraNonce2Length
	placeHolders := genPlaceHolders(placeHolderSize)

	payoutAddress, err := btcutil.DecodeAddress(payoutAddressStr, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, err
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	pkScript, err := txscript.PayToAddrScript(payoutAddress)

	tx.AddTxIn(&wire.TxIn{
		// Coinbase transactions have no inputs, so previous outpoint is
		// zero hash and max index.
		PreviousOutPoint: *wire.NewOutPoint(&chainhash.Hash{},
			wire.MaxPrevOutIndex),
		SignatureScript: addPlaceHolders(coinBaseSignatureScript, placeHolders),
		Sequence:        wire.MaxTxInSequenceNum,
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
	var transactions []*service.BlockTransactionPart

	if len(blockTemplateTransactions) > 0 {
		transactions = make([]*service.BlockTransactionPart, len(blockTemplateTransactions))
		txHashes = make([]*chainhash.Hash, len(blockTemplateTransactions))
		for i, tx := range blockTemplateTransactions {
			txHashes[i], _ = chainhash.NewHashFromStr(tx.Hash)
			transactions[i] = service.NewBlockTransactionPart()
			transactions[i].Data = tx.Data
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

func (c *BTCCoin) MakeBlock(header *service.BlockHeaderPart, base *service.BlockCoinBasePart, transactions []*service.BlockTransactionPart) (*service.Block, error) {
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

	coinBaseMsgTx := wire.NewMsgTx(wire.TxVersion)
	tx, err := hex.DecodeString(coinBaseBuf.String())
	if err != nil {
		return nil, err
	}

	err = coinBaseMsgTx.Deserialize(bytes.NewReader(tx))
	if err != nil {
		return nil, err
	}

	coinBaseTx := btcutil.NewTx(coinBaseMsgTx)
	blockTxns := []*btcutil.Tx{coinBaseTx}

	for _, transaction := range transactions {
		tx, err := decodeTransaction(transaction.Data)
		if err != nil {
			return nil, err
		}
		blockTxns = append(blockTxns, tx)
	}

	merkles := blockchain.BuildMerkleTreeStore(blockTxns, false)
	block.Header = wire.BlockHeader{
		Version:    blockVersion,
		PrevBlock:  prevHash,
		MerkleRoot: *merkles[len(merkles)-1],
		Timestamp:  timestamp,
		Bits:       bits,
		Nonce:      nonce,
	}

	blockHeaderHash := block.Header.BlockHash()
	for _, tx := range blockTxns {
		if err := block.AddTransaction(tx.MsgTx()); err != nil {
			return nil, err
		}
	}
	utilBlock := btcutil.NewBlock(&block)

	b := service.NewBlock()
	b.Hash = blockHeaderHash.String()
	b.Data = utilBlock.Hash().String()
	return b, nil
}

func (c *BTCCoin) SubmitBlock(data string) (bool, error) {
	var submitBlockParams [1]interface{}
	submitBlockParams[0] = data
	result, err := c.rpcClient.SubmitBlock(submitBlockParams)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (c *BTCCoin) IsSolveHash(hash string, targetDifficulty *big.Int) (bool, error) {
	hashPtr, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return false, err
	}

	if blockchain.HashToBig(hashPtr).Cmp(targetDifficulty) < 0 {
		return true, nil
	}
	return false, nil
}

func (c *BTCCoin) GetTargetDifficulty(bitsHex string) (*big.Int, error) {
	bits, err := decodeBits(bitsHex)
	if err != nil {
		return nil, err
	}
	return blockchain.CompactToBig(bits), nil
}

func (c *BTCCoin) CalculateShareComputePower(targetDifficulty *big.Int) (*big.Int, error) {
	//TODO check variable
	baseDifficultyBig := new(big.Int).SetUint64(BaseDifficulty)
	return new(big.Int).Mul(baseDifficultyBig, targetDifficulty), nil
}


func (c *BTCCoin) GetNewBlockHeight() (int, error) {
	miningInfo, err := c.rpcClient.GetMiningInfo()
	if err != nil {
		return 0, err
	}
	return miningInfo.Blocks + 1, nil
}