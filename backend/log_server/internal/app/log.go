package app

import (
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/himanhimao/lakepool/backend/log_server/internal/pkg/conf"
	pb "github.com/himanhimao/lakepool/backend/proto_log"
	"context"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"fmt"
	"strings"
	"strconv"
	"time"
	"errors"
	"encoding/json"
)

type LogServer struct {
	dbClient      client.Client
	config        *conf.LogConfig
	blockBPConfig client.BatchPointsConfig
	shareBPConfig client.BatchPointsConfig
}

func NewLogServer(client client.Client, conf *conf.LogConfig, blockBPConfig client.BatchPointsConfig, shareBPConfig client.BatchPointsConfig) *LogServer {
	return &LogServer{dbClient: client, config: conf, blockBPConfig: blockBPConfig, shareBPConfig: shareBPConfig}
}

func (s *LogServer) GetDBClient() client.Client {
	return s.dbClient
}

func (s *LogServer) GetConf() *conf.LogConfig {
	return s.config
}

func (s *LogServer) GetBlockBPConfig() client.BatchPointsConfig {
	return s.blockBPConfig
}

func (s *LogServer) GetShareBPConfig() client.BatchPointsConfig {
	return s.shareBPConfig
}

func (s *LogServer) QueryShareLogLatestTs(ctx context.Context, in *pb.QueryShareLogLatestTsRequest) (*pb.QueryShareLogLatestTsResponse, error) {
	if len(in.CoinType) == 0 {
		st := status.New(codes.InvalidArgument, "invalid coin type")
		return nil, st.Err()
	}

	if len(in.Tags.Hostname) == 0 {
		st := status.New(codes.InvalidArgument, "invalid tag host name")
		return nil, st.Err()
	}

	if in.Tags.Pid < 0 {
		st := status.New(codes.InvalidArgument, "invalid tag pid")
		return nil, st.Err()
	}

	var measurement string = fmt.Sprintf("%s_%s", s.config.MeasurementSharePrefix, strings.ToLower(in.CoinType))
	var command string = fmt.Sprintf("SELECT last(\"compute_power\") FROM %s WHERE \"host_name\"='%s' and \"pid\"='%d'",
		measurement, in.Tags.Hostname, in.Tags.Pid)

	query := client.NewQuery(command, s.shareBPConfig.Database, s.shareBPConfig.Precision)
	resp, err := s.dbClient.Query(query)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if resp.Err != "" {
		st := status.New(codes.Internal, errors.New(resp.Err).Error())
		return nil, st.Err()
	}

	result := resp.Results[0]

	if result.Err != "" {
		st := status.New(codes.Internal, errors.New(result.Err).Error())
		return nil, st.Err()
	}

	var latestTs int64
	fmt.Println("========================", result.Series, len(result.Series))
	if len(result.Series) > 0 {
		latestTs, _ = result.Series[0].Values[0][0].(json.Number).Int64()
	}

	return &pb.QueryShareLogLatestTsResponse{LatestTs: latestTs}, nil
}

func (s *LogServer) AddBlockLog(ctx context.Context, in *pb.AddBlockLogRequest) (*pb.AddBlockLogResponse, error) {
	if len(in.CoinType) == 0 {
		st := status.New(codes.InvalidArgument, "invalid coin type")
		return nil, st.Err()
	}

	measurement := fmt.Sprintf("%s_%s", s.config.MeasurementBlockPrefix, strings.ToLower(in.CoinType))
	fields := make(map[string]interface{})
	tags := make(map[string]string)
	tags["worker_name"] = in.Log.GetWorkerName()
	tags["server_ip"] = in.Log.GetServerIp()
	tags["client_ip"] = in.Log.GetClientIp()
	tags["user_name"] = in.Log.GetUserName()
	tags["user_agent"] = in.Log.GetUserAgent()
	tags["ext_name"] = in.Log.GetExtName()
	tags["host_name"] = in.Log.GetHostName()
	tags["pid"] = strconv.Itoa(int(in.Log.GetPid()))
	tags["height"] = strconv.Itoa(int(in.Log.GetHeight()))

	fields["hash"] = in.Log.Hash
	t := time.Unix(0, in.Ts)
	point, err := client.NewPoint(measurement, tags, fields, t)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	bp, _ := client.NewBatchPoints(s.GetBlockBPConfig())
	bp.AddPoint(point)

	if err := s.dbClient.Write(bp); err != nil {
		st := status.New(codes.Internal, err.Error())
		return &pb.AddBlockLogResponse{Result: false}, st.Err()
	}

	return &pb.AddBlockLogResponse{Result: true}, nil
}

func (s *LogServer) AddMinShareLogs(ctx context.Context, in *pb.AddMinShareLogsRequest) (*pb.AddMinShareLogsResponse, error) {
	if len(in.CoinType) == 0 {
		st := status.New(codes.InvalidArgument, "invalid coin type")
		return nil, st.Err()
	}

	measurement := fmt.Sprintf("%s_%s", s.config.MeasurementSharePrefix, strings.ToLower(in.CoinType))
	bp, _ := client.NewBatchPoints(s.GetShareBPConfig())
	for _, log := range in.Logs {
		fields := make(map[string]interface{})
		tags := make(map[string]string)
		tags["host_name"] = log.HostName
		tags["pid"] = strconv.Itoa(int(log.Pid))
		tags["server_ip"] = log.ServerIp
		tags["client_ip"] = log.ClientIp
		tags["user_name"] = log.UserName
		tags["ext_name"] = log.ExtName
		tags["worker_name"] = log.WorkerName
		tags["height"] = strconv.Itoa(int(log.Height))
		tags["user_agent"] = log.UserAgent
		if log.IsRight {
			tags["is_right"] = "1"
		} else {
			tags["is_right"] = "0"
		}

		fields["compute_power"] = log.ComputePower
		t := time.Unix(int64(log.Tm)*60, 0)
		point, err := client.NewPoint(measurement, tags, fields, t)
		if err != nil {
			st := status.New(codes.InvalidArgument, err.Error())
			return nil, st.Err()
		}
		bp.AddPoint(point)
	}

	if err := s.dbClient.Write(bp); err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	return &pb.AddMinShareLogsResponse{Result: true}, nil
}
