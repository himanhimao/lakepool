package main

import (
	pb "github.com/himanhimao/lakepool/backend/proto_stats"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
	"math/rand"
	"strings"
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

	workerNameSamples := []string{
		"himanhimao.001", "maomao.002", "gougou.003", "zhuzhu.004", "chichi.005",
	}

	serverIpSamples := []string{
		"11.11.11.11", "11.12.12.12", "11.13.13.13", "11.14.14.14", "11.15.15.15", "11.16.16.16", "11.17.17.17",
	}

	clientIpSamples := []string{
		"22.22.22.22", "22.12.12.12", "22.13.13.13", "22.14.14.14", "22.15.15.15", "22.16.16.16", "22.17.17.17",
	}

	userAgentSamples := []string{
		"aaa", "bbb", "ccc", "ddd", "eeee", "ffff", "gggg", "hhhh", "iiii",
	}

	hostNameSamples := []string{
		"localhost", "mycomputer", "mac-pro", "windows2008",
	}

	pidSamples := []int{1, 2, 3, 4, 5, 6}

	heightSamples := []int{564442, 564443, 564444, 564445, 564447, 564448}

	isRightSamples := []bool{true}

	sampleSize := 10000
	for i := 0; i < sampleSize; i++ {
		workerName := workerNameSamples[rand.Intn(len(workerNameSamples))]
		workerData := strings.Split(workerName, ".")

		addShareLogRequest := &pb.AddShareLogRequest{
			Ts:       time.Now().UnixNano(),
			CoinType: "BTC",
			Log: &pb.ShareLog{
				WorkerName:   workerName,
				ServerIp:     serverIpSamples[rand.Intn(len(serverIpSamples))],
				ClientIp:     clientIpSamples[rand.Intn(len(clientIpSamples))],
				UserName:     workerData[0],
				ExtName:      workerData[1],
				UserAgent:    userAgentSamples[rand.Intn(len(userAgentSamples))],
				HostName:     hostNameSamples[rand.Intn(len(hostNameSamples))],
				Pid:          int32(pidSamples[rand.Intn(len(pidSamples))]),
				Height:       int32(heightSamples[rand.Intn(len(heightSamples))]),
				ComputePower: float64(rand.Uint64() * 100),
				IsRight:      isRightSamples[rand.Intn(len(isRightSamples))],
			},
		}

		r, err := c.AddShareLog(ctx, addShareLogRequest)
		if err != nil {
			log.Fatalf("could not : %v", err)
		}
		log.Printf("%d-reult: %v", i, r.Result)
	}
}
