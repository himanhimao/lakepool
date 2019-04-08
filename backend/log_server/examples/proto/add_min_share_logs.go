package main

import (
	pb "github.com/himanhimao/lakepool/backend/proto_log"
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

	addMinShareLogsRequest := &pb.AddMinShareLogsRequest{
		CoinType: "BTC",
		Logs: []*pb.MinShareLog {
			{
				Tm: 2022332,
				WorkerName : "himan.001",
				ServerIp:  "127.0.0.1",
				ClientIp:  "192.168.1.1",
				UserName:  "himan",
				ExtName:   "001",
				UserAgent:  "test agent",
				HostName:   "xxxx",
				Pid:   1,
				Height: 22222,
				ComputePower: 300002,
				IsRight: true,
			},
			{
				Tm: 2022333,
				WorkerName : "himan.002",
				ServerIp:  "127.0.0.2",
				ClientIp:  "192.168.1.2",
				UserName:  "himan",
				ExtName:   "002",
				UserAgent:  "test xxx agent",
				HostName:   "xxxx",
				Pid:   1,
				Height: 222223,
				ComputePower: 300003,
				IsRight: true,
			},
		},
	}

	r, err := c.AddMinShareLogs(ctx, addMinShareLogsRequest)
	if err != nil {
		log.Fatalf("could not : %v", err)
	}
	log.Printf("result: %v", r.Result)
}
