package app

import (
	"github.com/himanhimao/lakepool/backend/stats_server/internal/pkg/worker"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_stats"
	"context"
	"github.com/himanhimao/lakepool/backend/stats_server/internal/pkg/conf"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/influxdata/influxdb1-client/v2"
	"fmt"
	"time"
	"strconv"
	"strings"
)

type StatsServer struct {
	worker *worker.StoreWorker
	config *conf.StatsConfig
}

func NewStatsServer(worker *worker.StoreWorker, config *conf.StatsConfig) *StatsServer {
	return &StatsServer{worker: worker, config: config}
}

func (s *StatsServer) AddShareLog(ctx context.Context, in *pb.AddShareLogRequest) (*pb.AddShareLogResponse, error) {
	if in.GetCoinType() != s.config.CoinType {
		st := status.New(codes.InvalidArgument, "Abnormal - invalid coin type")
		return nil, st.Err()
	}

	measurement := fmt.Sprintf("%s_%s", s.config.MeasurementSharePrefix, strings.ToLower(s.config.CoinType))
	fields := make(map[string]interface{})
	tags := make(map[string]string)
	tags["worker_name"] = in.Log.GetWorkerName()


	if in.Log.IsRight {
		tags["is_right"] = "1"
	} else {
		tags["is_right"] = "0"
	}

	fields["server_ip"] = in.Log.GetServerIp()
	fields["client_ip"] = in.Log.GetClientIp()
	fields["user_name"] = in.Log.GetUserName()
	fields["ext_name"] = in.Log.GetExtName()
	fields["user_agent"] = in.Log.GetUserAgent()
	fields["host_name"] = in.Log.GetHostName()
	fields["pid"] = strconv.Itoa(int(in.Log.GetPid()))
	fields["compute_power"] = in.Log.ComputePower
	fields["height"] = strconv.Itoa(int(in.Log.GetHeight()))

	t := time.Unix(0, in.Ts)
	sharePoint, err := client.NewPoint(measurement, tags, fields, t)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if err := s.worker.AsyncStoreSharePoint(sharePoint); err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}
	return &pb.AddShareLogResponse{Result: true}, nil
}

