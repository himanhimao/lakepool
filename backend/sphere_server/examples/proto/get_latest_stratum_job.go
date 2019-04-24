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

	r, err := c.GetLatestStratumJob(ctx, &pb.GetLatestStratumJobRequest{RegisterId:"1a1c5998"})
	if err != nil {
		log.Fatalf("could not : %v", err)
	}
	log.Printf("job: %v", r.Job)
}
