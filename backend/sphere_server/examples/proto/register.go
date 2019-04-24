package proto

import (
	pb "github.com/himanhimao/lakepool_proto/backend/proto_sphere"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:80"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSphereClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	//1K3kCJFD2PYF99t2eBSQ3stmdA1jrp7Nyf
	//r, err := c.GetBlockHeight(ctx, &pb.GetBlockHeightRequest{Hash: hash})
	r, err := c.Register(ctx, &pb.RegisterRequest{Config: &pb.StratumConfig{PoolTag: "/lakepool/", PayoutAddress: "1K3kCJFD2PYF99t2eBSQ3stmdA1jrp7Nyf", CoinType: "BTC"}, SysInfo: &pb.SysInfo{Hostname: "localhost", Pid: 10}})
	if err != nil {
		log.Fatalf("could not : %v", err)
	}
	log.Printf("register_id: %v", r.RegisterId)
}
