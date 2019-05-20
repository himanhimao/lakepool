package app

import (
	pb "github.com/himanhimao/lakepool_proto/backend/proto_log"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/himanhimao/lakepool/backend/stats_server/internal/pkg/conf"
	"sync"
	log "github.com/sirupsen/logrus"
	"time"
	"fmt"
	"encoding/json"
	"errors"
	"strconv"
	"context"
	"github.com/influxdata/influxdb1-client/models"
	"strings"
)

const (
	TIMEOUT = 10
)

type SyncTags struct {
	WorkerName string
	IsRight    int
}

type PushWorker struct {
	logGRPCClient pb.LogClient
	dbClient      client.Client
	conf          *conf.PusherConfig
	stopOnce      sync.Once
	runOnce       sync.Once
	stopC         chan struct{}
	syncTagsTimes map[SyncTags]int64
}

func NewPushWorker(logGRPCClient pb.LogClient, dbClient client.Client, conf *conf.PusherConfig) *PushWorker {
	return &PushWorker{logGRPCClient: logGRPCClient, dbClient: dbClient, conf: conf, stopC: make(chan struct{}),
		syncTagsTimes: make(map[SyncTags]int64)}
}

func (p *PushWorker) Run() {
	log.Infoln("pusher start...")
	p.runOnce.Do(func() {
		for {
			select {
			case <-p.stopC:
				return
			default:
				if err := p.sync(); err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Error("sync failed")
				} else {
					log.Debugln("sync success")
				}
			}
			time.Sleep(time.Second * p.conf.ReadInterval)
		}
	})
}

func (p *PushWorker) Stop() {
	p.stopOnce.Do(func() {
		log.Infoln("pusher stop...")
		p.stopC <- struct{}{}
	})
}

func (p *PushWorker) sync() error {
	result, err := p.getLocalGroupSeries()
	if err != nil {
		return err
	}

	for _, series := range result.Series {
		var mustUpdate bool
		var refTs int64
		var localCount int64

		workerName := series.Tags["worker_name"]
		isRight, _ := strconv.Atoi(series.Tags["is_right"])
		localLastTs, _ := series.Values[0][0].(json.Number).Int64()
		tags := SyncTags{WorkerName: workerName, IsRight: isRight}
		targetLastTs, err := p.getTargetLastTs(tags)
		if err != nil {
			return err
		}

		localCount, err = p.getLocalCount(tags, refTs)
		if err != nil {
			return err
		}
		
		if targetLastTs < localLastTs && localCount > 0 {
			mustUpdate = true
			refTs = targetLastTs
		} else {
			log.WithFields(log.Fields{
				"worker_name":    workerName,
				"is_right":       isRight,
				"target_last_ts": targetLastTs,
				"local_last_ts":  localLastTs,
				"local_count": localCount,
			}).Debugln("No need to update")
			continue
		}

		if mustUpdate {
			pageTotal := (localCount / int64(p.conf.PerSize)) + 1
			limit := p.conf.PerSize
			pageNum := int(0)

			var syncErr error
			for pageNum < int(pageTotal) {
				offset := pageNum * p.conf.PerSize
				row, err := p.getLocalSeries(tags, refTs, limit, offset)
				if err != nil {
					syncErr = errors.New(fmt.Sprintf("workerName %s, isRight $d limit %d, offset %d:%s", tags.WorkerName, tags.IsRight, limit, offset, err.Error()))
					break
				}
				rowsLen := len(row.Values)
				if rowsLen > 0 {
					pbMinShareLogs := make([]*pb.MinShareLog, rowsLen)
					for index, value := range row.Values {
						pbMinShareLog := new(pb.MinShareLog)
						for i, column := range row.Columns {
							if column == "time" {
								tm, _ := value[i].(json.Number).Int64()
								pbMinShareLog.Tm = int32(tm)
							} else if column == "client_ip" {
								pbMinShareLog.ClientIp = value[i].(string)
							} else if column == "compute_power" {
								computePower, _ := value[i].(json.Number).Float64()
								pbMinShareLog.ComputePower = computePower
							} else if column == "ext_name" {
								pbMinShareLog.ExtName = value[i].(string)
							} else if column == "height" {
								height, _ := strconv.Atoi(value[i].(string))
								pbMinShareLog.Height = int32(height)
							} else if column == "host_name" {
								pbMinShareLog.HostName = value[i].(string)
							} else if column == "is_right" {
								isRight, _ := strconv.ParseBool(value[i].(string))
								pbMinShareLog.IsRight = isRight
							} else if column == "pid" {
								pid, _ := strconv.Atoi(value[i].(string))
								pbMinShareLog.Pid = int32(pid)
							} else if column == "server_ip" {
								pbMinShareLog.ServerIp = value[i].(string)
							} else if column == "user_agent" {
								pbMinShareLog.UserAgent = value[i].(string)
							} else if column == "user_name" {
								pbMinShareLog.UserName = value[i].(string)
							} else if column == "worker_name" {
								pbMinShareLog.WorkerName = value[i].(string)
							}
						}
						pbMinShareLogs[index] = pbMinShareLog
					}

					pbAddMinShareLogsRequest := &pb.AddMinShareLogsRequest{
						Logs:     pbMinShareLogs,
						CoinType: p.conf.CoinType,
					}

					r, err := p.logGRPCClient.AddMinShareLogs(context.Background(), pbAddMinShareLogsRequest)
					if err != nil {
						log.Fatalf("could not : %v", err)
					}

					if r.Result {
						p.syncTagsTimes[tags] = int64(pbMinShareLogs[len(pbMinShareLogs)-1].Tm)
						log.WithFields(log.Fields{
							"worker_name": tags.WorkerName,
							"is_right":    tags.IsRight,
						}).Info("add min share logs success")
					} else {
						log.WithFields(log.Fields{
							"worker_name": tags.WorkerName,
							"is_right":    tags.IsRight,
						}).Error("add min share logs failed")
					}
				}
				pageNum++
			}

			if syncErr != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PushWorker) getTargetLastTs(tags SyncTags) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*TIMEOUT)
	defer cancel()

	var isRight bool
	if tags.IsRight > 0 {
		isRight = true
	} else {
		isRight = false
	}

	targetResp, err := p.logGRPCClient.QueryShareLogLatestTs(ctx, &pb.QueryShareLogLatestTsRequest{
		Tags: &pb.QueryShareLogTags{
			WorkerName: tags.WorkerName,
			IsRight:    isRight,
		},
		CoinType: p.conf.CoinType,
	})

	if err != nil {
		return 0, err
	}
	return targetResp.LatestTs, nil
}

func (p *PushWorker) getLocalGroupSeries() (client.Result, error) {
	var query client.Query
	var command = fmt.Sprintf("SELECT last(\"compute_power\") FROM %s GROUP BY \"worker_name\",\"is_right\"",
		p.formatMeasurement())

	if len(p.conf.RetentionStrategy) > 0 {
		query = client.NewQueryWithRP(command, p.conf.Database, p.conf.RetentionStrategy, p.conf.Precision)
	} else {
		query = client.NewQuery(command, p.conf.Database, p.conf.Precision)
	}
	var result client.Result

	resp, err := p.dbClient.Query(query)
	if err != nil {
		return result, err
	}

	if resp.Err != "" {
		return result, errors.New(resp.Err)
	}

	if result.Err != "" {
		return result, errors.New(result.Err)
	}

	return resp.Results[0], nil
}

func (p *PushWorker) getLocalCount(tags SyncTags, ts int64) (int64, error) {
	var err error
	var count int64
	defer func() {
		recover()
	}()

	var query client.Query
	var command string
	if ts > 0 {
		command = fmt.Sprintf("SELECT count(*) FROM %s WHERE \"worker_name\"='%s' and \"is_right\"='%d'"+
			" and \"time\">%d%s", p.formatMeasurement(), tags.WorkerName, tags.IsRight, ts, p.conf.Precision)
	} else {
		command = fmt.Sprintf("SELECT count(*) FROM %s WHERE \"worker_name\"='%s' and \"is_right\"='%d'",
			p.formatMeasurement(), tags.WorkerName, tags.IsRight)
	}

	if len(p.conf.RetentionStrategy) > 0 {
		query = client.NewQueryWithRP(command, p.conf.Database, p.conf.RetentionStrategy, p.conf.Precision)
	} else {
		query = client.NewQuery(command, p.conf.Database, p.conf.Precision)
	}

	resp, err := p.dbClient.Query(query)
	if err != nil {
		return count, err
	}

	if resp.Err != "" {
		return count, errors.New(resp.Err)
	}

	count, err = resp.Results[0].Series[0].Values[0][1].(json.Number).Int64()
	return count, err

}

func (p *PushWorker) getLocalSeries(tags SyncTags, ts int64, limit int, offset int) (*models.Row, error) {
	var err error
	var query client.Query
	var command string

	command = fmt.Sprintf("SELECT * FROM %s WHERE \"worker_name\"='%s' and \"is_right\"='%d'"+
		" and \"time\">%d%s LIMIT %d OFFSET %d", p.formatMeasurement(), tags.WorkerName, tags.IsRight, ts,
		p.conf.Precision, limit, offset)

	if len(p.conf.RetentionStrategy) > 0 {
		query = client.NewQueryWithRP(command, p.conf.Database, p.conf.RetentionStrategy, p.conf.Precision)
	} else {
		query = client.NewQuery(command, p.conf.Database, p.conf.Precision)
	}

	resp, err := p.dbClient.Query(query)
	if err != nil {
		return nil, err
	}

	if resp.Err != "" {
		return nil, errors.New(resp.Err)
	}

	return &resp.Results[0].Series[0], nil
}

func (p *PushWorker) formatMeasurement() string {
	return fmt.Sprintf("%s_%s_%s", p.conf.MeasurementPrefix, strings.ToLower(p.conf.CoinType), p.conf.MeasurementSuffix)
}
