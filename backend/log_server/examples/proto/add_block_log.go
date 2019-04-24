package main

import (
	pb "github.com/himanhimao/lakepool_proto/backend/proto_log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:8082"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLogClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	addBlockLogRequest := &pb.AddBlockLogRequest{
		Ts:       time.Now().UnixNano(),
		CoinType: "BTC",
		Log: &pb.BlockLog{
			WorkerName: "himanhimao.001",
			ServerIp:   "11.11.11.11",
			ClientIp:   "22.22.22.22",
			UserName:   "himanhimao",
			ExtName:    "001",
			UserAgent:  "abc",
			HostName:   "localhost",
			Pid:        1,
			Height:     564442,
			Hash:    "000000000000000000283507be94843e450143c381c6d7e6b141c439c172472c",
		},
	}

	r, err := c.AddBlockLog(ctx, addBlockLogRequest)
	if err != nil {
		log.Fatalf("could not : %v", err)
	}
	log.Printf("result: %v", r.Result)
}
