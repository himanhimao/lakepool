package main

import (
	pb "github.com/himanhimao/lakepool/backend/proto_stats"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:8080"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStatsClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	addShareLogRequest := &pb.AddShareLogRequest{
		Ts:       time.Now().UnixNano(),
		CoinType: "BTC",
		Log: &pb.ShareLog{
			WorkerName: "himanhimao.001",
			ServerIp:   "11.11.11.11",
			ClientIp:   "22.22.22.22",
			UserName:   "himanhimao",
			ExtName:    "001",
			UserAgent:  "abc",
			HostName:   "localhost",
			Pid:        1,
			Height:     564442,
			ComputePower:    16384 * 60,
			IsRight:    true,
		},
	}

	r, err := c.AddShareLog(ctx, addShareLogRequest)
	if err != nil {
		log.Fatalf("could not : %v", err)
	}
	log.Printf("reult: %v", r.Result)
}
