package btc

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"bytes"
	"io/ioutil"
)

const (
	rpcClientTimeout       int = 30
	VersionV1                  = "1.0"
	VersionV2                  = "2.0"
	methodGetBlockTemplate     = "getblocktemplate"
	methodSubmitBlock          = "submitblock"
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

// call prepare & exec the request
func (c *RpcClient) call1(method string, params interface{}, version string) (rr RpcResponse, err error) {
	data := `{"result":{"capabilities":["proposal"],"version":536870912,"rules":["csv","!segwit"],"vbavailable":{},"vbrequired":0,"previousblockhash":"0000000000000047e5bda122407654b25d52e0f3eeb00c152f631f70e9803772","transactions":[{"data":"0100000002449f651247d5c09d3020c30616cb1807c268e2c2346d1de28442b89ef34c976d000000006a47304402203eae3868946a312ba712f9c9a259738fee6e3163b05d206e0f5b6c7980161756022017827f248432f7313769f120fb3b7a65137bf93496a1ae7d6a775879fbdfb8cd0121027d7b71dab3bb16582c97fc0ccedeacd8f75ebee62fa9c388290294ee3bc3e935feffffffcbc82a21497f8db8d57d054fefea52aba502a074ed984efc81ec2ef211194aa6010000006a47304402207f5462295e52fb4213f1e63802d8fe9ec020ac8b760535800564694ea87566a802205ee01096fc9268eac483136ce082506ac951a7dbc9e4ae24dca07ca2a1fdf2f30121023b86e60ef66fe8ace403a0d77d27c80ba9ba5404ee796c47c03c73748e59d125feffffff0286c35b00000000001976a914ab29f668d284fd2d65cec5f098432c4ece01055488ac8093dc14000000001976a914ac19d3fd17710e6b9a331022fe92c693fdf6659588ac8dd70f00","txid":"c284853b65e7887c5fd9b635a932e2e0594d19849b22914a8e6fb180fea0954f","hash":"c284853b65e7887c5fd9b635a932e2e0594d19849b22914a8e6fb180fea0954f","depends":[],"fee":37400,"sigops":8,"weight":1488},{"data":"0100000001043f5e73755b5c6919b4e361f4cae84c8805452de3df265a6e2d3d71cbcb385501000000da0047304402202b14552521cd689556d2e44d914caf2195da37b80de4f8cd0fad9adf7ef768ef022026fcddd992f447c39c48c3ce50c5960e2f086ebad455159ffc3e36a5624af2f501483045022100f2b893e495f41b22cd83df6908c2fa4f917fd7bce9f8da14e6ab362042e11f7d022075bc2451e1cf2ae2daec0f109a3aceb6558418863070f5e84c945262018503240147522102632178d046673c9729d828cfee388e121f497707f810c131e0d3fc0fe0bd66d62103a0951ec7d3a9da9de171617026442fcd30f34d66100fab539853b43f508787d452aeffffffff0240420f000000000017a9143e9a6b79be836762c8ef591cf16b76af1327ced58790dfdf8c0000000017a9148ce5408cfeaddb7ccb2545ded41ef478109454848700000000","txid":"28b1a5c2f0bb667aea38e760b6d55163abc9be9f1f830d9969edfab902d17a0f","hash":"28b1a5c2f0bb667aea38e760b6d55163abc9be9f1f830d9969edfab902d17a0f","depends":[],"fee":20000,"sigops":8,"weight":1332},{"data":"01000000013faf73481d6b96c2385b9a4300f8974b1b30c34be30000c7dcef11f68662de4501000000db00483045022100f9881f4c867b5545f6d7a730ae26f598107171d0f68b860bd973dbb855e073a002207b511ead1f8be8a55c542ce5d7e91acfb697c7fa2acd2f322b47f177875bffc901483045022100a37aa9998b9867633ab6484ad08b299de738a86ae997133d827717e7ed73d953022011e3f99d1bd1856f6a7dc0bf611de6d1b2efb60c14fc5931ba09da01558757f60147522102632178d046673c9729d828cfee388e121f497707f810c131e0d3fc0fe0bd66d62103a0951ec7d3a9da9de171617026442fcd30f34d66100fab539853b43f508787d452aeffffffff0240420f000000000017a9148d57003ecbaa310a365f8422602cc507a702197e87806868a90000000017a9148ce5408cfeaddb7ccb2545ded41ef478109454848700000000","txid":"67878210e268d87b4e6587db8c6e367457cea04820f33f01d626adbe5619b3dd","hash":"67878210e268d87b4e6587db8c6e367457cea04820f33f01d626adbe5619b3dd","depends":[],"fee":20000,"sigops":8,"weight":1336}],"coinbaseaux":{"flags":""},"coinbasevalue":319367518,"longpollid":"0000000000000047e5bda122407654b25d52e0f3eeb00c152f631f70e9803772604597","target":"0000000000001714480000000000000000000000000000000000000000000000","mintime":1480831053,"mutable":["time","transactions","prevblock"],"noncerange":"00000000ffffffff","sigoplimit":80000,"sizelimit":4000000,"weightlimit":4000000,"curtime":1480834892,"bits":"17306835","height":1038222,"default_witness_commitment":"6a24aa21a9ed842a6d6672504c2b7abb796fdd7cfbd7262977b71b945452e17fbac69ed22bf8"}}`
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
