package app

import (
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/himanhimao/lakepool/backend/log_server/internal/pkg/conf"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_log"
	"context"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"fmt"
	"strings"
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

func (s *LogServer) QueryShareLogLatestTs(ctx context.Context, in *pb.QueryShareLogLatestTsRequest) (*pb.QueryShareLogLatestTsResponse, error) {
	if len(in.CoinType) == 0 {
		st := status.New(codes.InvalidArgument, "invalid coin type")
		return nil, st.Err()
	}

	if len(in.Tags.WorkerName) == 0 {
		st := status.New(codes.InvalidArgument, "invalid tag worker name")
		return nil, st.Err()
	}

	var isRight int
	if in.Tags.IsRight {
		isRight = 1
	} else {
		isRight = 0
	}

	var measurement string = fmt.Sprintf("%s_%s", s.config.MeasurementSharePrefix, strings.ToLower(in.CoinType))
	var command string = fmt.Sprintf("SELECT last(\"compute_power\") FROM %s WHERE \"host_name\"='%s' and \"is_right\"='%d'", measurement,
		in.Tags.WorkerName, isRight)

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
	tags["hash"] = in.Log.Hash
	fields["server_ip"] = in.Log.GetServerIp()
	fields["client_ip"] = in.Log.GetClientIp()
	fields["user_name"] = in.Log.GetUserName()
	fields["user_agent"] = in.Log.GetUserAgent()
	fields["ext_name"] = in.Log.GetExtName()
	fields["host_name"] = in.Log.GetHostName()
	fields["pid"] = int(in.Log.GetPid())
	fields["height"] = int(in.Log.GetHeight())
	t := time.Unix(0, in.Ts)
	point, err := client.NewPoint(measurement, tags, fields, t)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	bp, _ := client.NewBatchPoints(s.blockBPConfig)
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
	bp, _ := client.NewBatchPoints(s.shareBPConfig)
	for _, log := range in.Logs {
		fields := make(map[string]interface{})
		tags := make(map[string]string)
		tags["worker_name"] = log.WorkerName
		if log.IsRight {
			tags["is_right"] = "1"
		} else {
			tags["is_right"] = "0"
		}

		fields["pid"] = int(log.Pid)
		fields["server_ip"] = log.ServerIp
		fields["client_ip"] = log.ClientIp
		fields["user_name"] = log.UserName
		fields["ext_name"] = log.ExtName
		fields["host_name"] = log.HostName
		fields["height"] = int(log.Height)
		fields["user_agent"] = log.UserAgent
		fields["host_name"] = log.HostName
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
