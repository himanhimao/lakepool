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

	queryShareLogLatestTsRequest := &pb.QueryShareLogLatestTsRequest{
		Tags: &pb.QueryShareLogTags{
			Pid:1,
			Hostname: "localhost",
		},
		MeasurementName: "btc_1min",
	}

	r, err := c.QueryShareLogLatestTs(ctx, queryShareLogLatestTsRequest)
	if err != nil {
		log.Fatalf("could not : %v", err)
	}
	log.Printf("result: %v", r.LatestTs)
}
