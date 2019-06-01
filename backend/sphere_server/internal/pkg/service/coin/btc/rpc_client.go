package btc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	rpcClientTimeout       int = 30
	VersionV1                  = "1.0"
	VersionV2                  = "2.0"
	methodGetBlockTemplate     = "getblocktemplate"
	methodSubmitBlock          = "submitblock"
	methodGetMiningInfo        = "getmininginfo"
)

// A rpcClient represents a JSON RPC client (over HTTP(s)).
type RpcClient struct {
	serverAddr string
	user       string
	passwd     string
	httpClient *http.Client
}

// rpcRequest represent a RCP request
type rpcRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int64       `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
}

type BlockTemplateTransaction struct {
	Data    string  `json:"data"`
	TxId    string  `json:"txid"`
	Hash    string  `json:"hash"`
	Depends []int32 `json:"depends"`
	Fee     int32   `json:"int32"`
	SigOps  int32   `json:"sigops"`
	Weight  int32   `json:"weight"`
}

type BlockTemplate struct {
	Capabilities             []string                    `json:"capabilities"`
	Version                  int32                       `json:"version"`
	Rules                    []string                    `json:"rules"`
	VBAvailable              interface{}                 `json:"vbavailable"`
	VBRequired               int32                       `json:"vbrequired"`
	PreviousBlockHash        string                      `json:"previousblockhash"`
	Transactions             []*BlockTemplateTransaction `json:"transactions"`
	CoinBaseAux              interface{}                 `json:"coinbaseaux"`
	CoinBaseValue            int64                       `json:"coinbasevalue"`
	LongPollId               string                      `json:"langpollid"`
	Target                   string                      `json:"target"`
	MinTime                  int32                       `json:"mintime"`
	Mutable                  []string                    `json:"mutable"`
	NonceRange               string                      `json:"noncerange"`
	SigOpLimit               int32                       `json:"sigoplimit"`
	SizeLimit                int32                       `json:"sizelimit"`
	WeightLimit              int32                       `json:"weightlimit"`
	CurTime                  int32                       `json:"curtime"`
	Bits                     string                      `json:"bits"`
	Height                   int32                       `json:"height"`
	DefaultWitnessCommitment string                      `json:"default_witness_commitment"`
}

type MiningInfo struct {
	Blocks             int     `json:"blocks"`
	CurrentBlockWeight int     `json:"currentblockweight"`
	CurrentBlockTx     int     `json:"currentblocktx"`
	Difficulty         float64 `json:"difficulty"`
	NetworkHashPs      float64 `json:"networkhashps"`
	PooledTx           int     `json:"pooledtx"`
	Chain              string  `json:"chain"`
	Warnings           string  `json:"warnings"`
}

// rpcError represents a RCP error
/*type rpcError struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}*/

type RpcResponse struct {
	Id     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    interface{}     `json:"error"`
}

func NewClient(host string, port int, user, passwd string, useSSL bool) (c *RpcClient) {
	//TODO validate
	var serverAddr string
	var httpClient *http.Client
	if useSSL {
		serverAddr = "https://"
		t := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: t}
	} else {
		serverAddr = "http://"
		httpClient = &http.Client{}
	}
	c = &RpcClient{serverAddr: fmt.Sprintf("%s%s:%d", serverAddr, host, port), user: user, passwd: passwd, httpClient: httpClient}
	return
}

// doTimeoutRequest process a HTTP request with timeout
func (c *RpcClient) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("Timeout reading data from server")
	}
}

// call prepare & exec the request
func (c *RpcClient) call(method string, params interface{}, version string) (rr RpcResponse, err error) {
	connectTimer := time.NewTimer(time.Duration(rpcClientTimeout) * time.Second)
	rpcR := rpcRequest{method, params, time.Now().UnixNano(), version}
	payloadBuffer := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(payloadBuffer)
	err = jsonEncoder.Encode(rpcR)
	if err != nil {
		return
	}
	//fmt.Println(payloadBuffer.String())
	req, err := http.NewRequest("POST", c.serverAddr, payloadBuffer)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	// Auth ?
	if len(c.user) > 0 || len(c.passwd) > 0 {
		req.SetBasicAuth(c.user, c.passwd)
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("HTTP error: " + resp.Status)
		return
	}
	err = json.Unmarshal([]byte(data), &rr)
	return
}

func (c *RpcClient) GetBlockTemplate(params interface{}) (*BlockTemplate, error) {
	resp, err := c.call(methodGetBlockTemplate, params, VersionV1)
	if err != nil {
		return nil, err
	}
	var blockTemplate BlockTemplate

	if err := json.Unmarshal(resp.Result, &blockTemplate); err != nil {
		return nil, err
	}
	return &blockTemplate, nil
}

func (c *RpcClient) SubmitBlock(params interface{}) (bool, error) {
	resp, err := c.call(methodSubmitBlock, params, VersionV1)
	if err != nil {
		return false, err
	}

	var result string
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return false, err
	}

	if len(result) > 0 {
		return false, errors.New(result)
	}
	return true, nil
}

func (c *RpcClient) GetMiningInfo() (*MiningInfo, error) {
	resp, err := c.call(methodGetMiningInfo, nil, VersionV1)
	if err != nil {
		return  nil, err
	}

	var miningInfo MiningInfo
	if err := json.Unmarshal(resp.Result, &miningInfo); err != nil {
		return nil, err
	}

	return &miningInfo, nil
}



