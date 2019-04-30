package cache

import (
	"github.com/gomodule/redigo/redis"
	"errors"
	"strconv"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/service"
	"encoding/json"
)

var (
	ErrorNoInitialization = errors.New("not initialization")
)

type RedisCache struct {
	pool *redis.Pool
}

func NewRedisCache() *RedisCache {
	return &RedisCache{}
}


func (c *RedisCache) SetClient(client *redis.Pool) *RedisCache {
	c.pool = client
	return c
}

func (c *RedisCache) SetRegisterContext(key service.RegisterKey, r *service.Register) error {
	client := c.pool.Get()
	defer client.Close()
	if client == nil {
		return ErrorNoInitialization
	}

	extraNonce1Length := r.ExtraNonce1Length
	extraNonce2Length := r.ExtraNonce2Length

	if _, err := client.Do("hmset", key, service.KeyRegPayoutAddress,
		r.PayoutAddress, service.KeyRegPoolTag, r.PoolTag, service.KeyRegCoinType,
		r.CoinType, service.KeyRegUsedTestNet, strconv.FormatBool(r.UsedTestNet),
		service.KeyRegExtraNonce1Length, strconv.Itoa(extraNonce1Length),
		service.KeyRegExtraNonce2Length, strconv.Itoa(extraNonce2Length),
	); err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) GetRegisterContext(key service.RegisterKey) (*service.Register, error) {
	client := c.pool.Get()
	defer client.Close()
	if client == nil {
		return nil, ErrorNoInitialization
	}
	values, err := redis.Values(client.Do("hmget", key, service.KeyRegPayoutAddress, service.KeyRegPoolTag,
		service.KeyRegCoinType, service.KeyRegUsedTestNet, service.KeyRegExtraNonce1Length, service.KeyRegExtraNonce2Length))
	if err != nil {
		return nil, err
	}

	r := service.NewRegister()
	if values[0] != nil {
		r.PayoutAddress = string(values[0].([]byte))
	}

	if values[1] != nil {
		r.PoolTag  = string(values[1].([]byte))
	}

	if values[2] != nil {
		r.CoinType = string(values[2].([]byte))
	}

	if values[3] != nil {
		res, _ := strconv.ParseBool(string(values[3].([]byte)))
		r.UsedTestNet = res
	}

	if values[4] != nil {
		res, _ := strconv.Atoi(string(values[4].([]byte)))
		r.ExtraNonce1Length = res
	}

	if values[5] != nil {
		res, _ := strconv.Atoi(string(values[5].([]byte)))
		r.ExtraNonce2Length = res
	}
	return r, nil
}

func (c *RedisCache) DelRegisterContext(key service.RegisterKey) error {
	client := c.pool.Get()
	defer client.Close()
	if client == nil {
		return ErrorNoInitialization
	}

	_, err := client.Do("del", key)
	return err
}

func (c *RedisCache) SetBlockTransactions(key service.JobKey, expireTs int, transactions []*service.BlockTransactionPart) error {
	client := c.pool.Get()
	defer client.Close()
	if client == nil {
		return ErrorNoInitialization
	}

	transactionsBytes, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	if _, err := client.Do("set", key, transactionsBytes); err != nil {
		return err
	}

	if _, err := client.Do("expire", key, expireTs); err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) GetBlockTransactions(key service.JobKey) ([]*service.BlockTransactionPart, error) {
	client := c.pool.Get()
	defer client.Close()
	if client == nil {
		return nil, ErrorNoInitialization
	}

	transactionsBytes, err := redis.Bytes(client.Do("get", key))
	if err != nil {
		return nil, err
	}

	var transactions = make([]*service.BlockTransactionPart, 0)
	err = json.Unmarshal(transactionsBytes, &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (c *RedisCache) SetShareHash(key service.ShareKey, hash string, ) error {
	client := c.pool.Get()
	defer client.Close()

	if client == nil {
		return ErrorNoInitialization
	}

	if _, err := client.Do("hset", key, hash); err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) ExistShareHash(key service.ShareKey, hash string) (bool, error) {
	client := c.pool.Get()
	defer client.Close()

	if client == nil {
		return false, ErrorNoInitialization
	}

	ok, err := redis.Bool(client.Do("hexists", key, hash))
	return ok, err
}

func (c *RedisCache) ClearShareHistory(key service.ShareKey) error {
	client := c.pool.Get()
	defer client.Close()

	if client == nil {
		return ErrorNoInitialization
	}

	_, err := client.Do("del", key)
	return err
}
