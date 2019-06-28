package proto

import (
	"encoding/json"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/util"
	"reflect"
)

const (
	MsgFullNameJsonRpcREQ                               = "proto.JSONRpcREQ"
	MsgFullNameJsonRpcResp                              = "proto.JSONRpcRESP"
	MsgFullNameJSONRpcSetDifficultyRESP                 = "proto.JSONRpcSetDifficultyRESP"
	MsgFullNameJSONRpcNotifyRESP                        = "proto.JSONRpcNotifyRESP"
	MsgFullNameJobNotify                                = "proto.jobNotify"
	MsgFullNameJobDifficultyCheck                       = "proto.jobDifficultyCheck"
	MsgFullNameJobHashCheck                             = "proto.jobHashCheck"
	MiningRecvMethodSubscribe                           = "mining.subscribe"
	MiningRecvMethodExtranonceSubscribe                 = "mining.extranonce.subscribe"
	MiningRecvMethodAuthorize                           = "mining.authorize"
	MiningRecvMethodSubmit                              = "mining.submit"
	MiningRespMethodSetDifficulty                       = "mining.set_difficulty"
	MiningRespMethodNotify                              = "mining.notify"
	ConnAcceptd                                         = "conn.accepted"
	ConnClosed                                          = "conn.closed"
	ConnUnknown                                         = "conn.unknown"
	JobMethodNotify                                     = "job.notify"
	JobMethodDifficultyCheck                            = "job.difficulty.check"
	JobMethodHeightCheck                                = "job.height.check"
	AuthorizePass                       AuthorizeResult = true
	AuthorizeRefuse                     AuthorizeResult = false
	SubmitPass                          SubmitResult    = true
	SubmitRefuse                        SubmitResult    = false
)

type AuthorizeResult bool
type SubmitResult bool

type Loader interface {
	Load(j *JSONRpcREQ)
}

type JSONRpcREQ struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params []interface{}    `json:"params"`
}

type JSONRpcRESP struct {
	Id     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  []interface{}    `json:"error"`
}

type JSONRpcSetDifficultyRESP struct {
	Method string   `json:"method"`
	Params []uint64 `json:"params"`
}

type JSONRpcNotifyRESP struct {
	Id     interface{}   `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type SubscribeREQ struct {
	Id        *json.RawMessage
	UserAgent string
	Session   string
}

type AuthorizeREQ struct {
	Id         *json.RawMessage
	WorkerName string
	Password   string
}

type SubmitREQ struct {
	Id          *json.RawMessage
	WorkerName  string
	JobId       string
	ExtraNonce2 string
	NTime       string
	Nonce       string
}

type JobNotify struct {
}

type JobDifficultyCheck struct {
}

type JobHeightCheck struct {
}

func (o *JSONRpcREQ) String() string               { return fmt.Sprintf("%+v", *o) }
func (o *JSONRpcRESP) String() string              { return fmt.Sprintf("%+v", *o) }
func (o *JSONRpcSetDifficultyRESP) String() string { return fmt.Sprintf("%+v", *o) }
func (o *SubscribeREQ) String() string             { return fmt.Sprintf("%+v", *o) }
func (o *AuthorizeREQ) String() string             { return fmt.Sprintf("%+v", *o) }
func (o *JSONRpcNotifyRESP) String() string        { return fmt.Sprintf("%+v", *o) }
func (o *JobNotify) String() string                { return fmt.Sprintf("%+v", *o) }
func (o *JobDifficultyCheck) String() string       { return fmt.Sprintf("%+v", *o) }
func (o *JobHeightCheck) String() string           { return fmt.Sprintf("%+v", *o) }

func NewSubscribeREQ() *SubscribeREQ {
	return &SubscribeREQ{}
}

func NewAuthorizeREQ() *AuthorizeREQ {
	return &AuthorizeREQ{}
}

func NewSubmitREQ() *SubmitREQ {
	return &SubmitREQ{}
}

func NewErrOtherUnknownRESP(id *json.RawMessage) *JSONRpcRESP {
	return NewErrResp(id, ErrOtherUnknown)
}

func NewErrNotSubscribedRESP(id *json.RawMessage) *JSONRpcRESP {
	return NewErrResp(id, ErrNotSubscribed)
}

func NewErrUnauthorizedWorkerRESP(id *json.RawMessage) *JSONRpcRESP {
	return NewErrResp(id, ErrUnauthorizedWorker)
}

func NewErrJobNotFoundRESP(id *json.RawMessage) *JSONRpcRESP {
	return NewErrResp(id, ErrJobNotFound)
}

func NewErrResp(id *json.RawMessage, errType ErrType) *JSONRpcRESP {
	resp := new(JSONRpcRESP)
	resp.Id = id
	resp.Error = make([]interface{}, 2)
	resp.Error[0] = errType
	resp.Error[1] = errType.ErrMsg()
	return resp
}

func NewSubscribeRESP(id *json.RawMessage, notify interface{}, difficulty interface{}, extraNonce1 string, extraNonce2Length int) *JSONRpcRESP {
	resp := new(JSONRpcRESP)
	resp.Id = id
	resp.Error = nil
	result := make([]interface{}, 3)
	result[0] = [2][2]interface{}{{MiningRespMethodSetDifficulty, difficulty}, {MiningRespMethodNotify, notify}}
	result[1] = extraNonce1
	result[2] = extraNonce2Length
	resp.Result = result
	return resp
}

func NewAuthorizeRESP(id *json.RawMessage, result AuthorizeResult) *JSONRpcRESP {
	resp := new(JSONRpcRESP)
	resp.Id = id
	resp.Result = result
	resp.Error = nil
	return resp
}

func NewSubmitRESP(id *json.RawMessage, result SubmitResult) *JSONRpcRESP {
	resp := new(JSONRpcRESP)
	resp.Id = id
	resp.Result = result
	resp.Error = nil
	return resp
}

func NewJSONRpcSetDifficultyRESP(difficulty uint64) *JSONRpcSetDifficultyRESP {
	resp := new(JSONRpcSetDifficultyRESP)
	resp.Method = MiningRespMethodSetDifficulty
	params := make([]uint64, 1)
	params[0] = difficulty
	resp.Params = params
	return resp
}

func NewJSONRpcNotifyRESP(params []interface{}) *JSONRpcNotifyRESP {
	resp := new(JSONRpcNotifyRESP)
	resp.Id = nil
	resp.Params = params
	resp.Method = MiningRespMethodNotify
	return resp
}

func (o *SubscribeREQ) Load(j *JSONRpcREQ) {
	o.Id = j.Id
	for i, val := range j.Params {
		if i > 1 {
			break
		}

		if i == 0 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.UserAgent = val.(string)
			}
		}

		if i == 1 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.Session = val.(string)
			}
		}
	}
}

func (o *AuthorizeREQ) Load(j *JSONRpcREQ) {
	o.Id = j.Id
	for i, val := range j.Params {
		if i > 1 {
			break
		}

		if i == 0 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.WorkerName = val.(string)
			}
		}

		if i == 1 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.Password = val.(string)
			}
		}
	}
}

func (o *SubmitREQ) Load(j *JSONRpcREQ) {
	o.Id = j.Id
	for i, val := range j.Params {
		if i > 4 {
			break
		}

		if i == 0 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.WorkerName = val.(string)
			}
		}

		if i == 1 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.JobId = val.(string)
			}
		}

		if i == 2 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.ExtraNonce2 = val.(string)
			}
		}

		if i == 3 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.NTime = val.(string)
			}
		}

		if i == 4 {
			if reflect.TypeOf(val).Kind() == reflect.String {
				o.Nonce = val.(string)
			}
		}
	}
}

// 将消息注册到系统
func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*JSONRpcREQ)(nil)).Elem(),
		ID:    int(util.StringHash(MsgFullNameJsonRpcREQ)),
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*JSONRpcRESP)(nil)).Elem(),
		ID:    int(util.StringHash(MsgFullNameJsonRpcResp)),
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*JSONRpcSetDifficultyRESP)(nil)).Elem(),
		ID:    int(util.StringHash(MsgFullNameJSONRpcSetDifficultyRESP)),
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*JSONRpcNotifyRESP)(nil)).Elem(),
		ID:    int(util.StringHash(MsgFullNameJSONRpcNotifyRESP)),
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*JobNotify)(nil)).Elem(),
		ID:    int(util.StringHash(MsgFullNameJobNotify)),
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*JobDifficultyCheck)(nil)).Elem(),
		ID:    int(util.StringHash(MsgFullNameJobDifficultyCheck)),
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*JobHeightCheck)(nil)).Elem(),
		ID:    int(util.StringHash(MsgFullNameJobHashCheck)),
	})
}
