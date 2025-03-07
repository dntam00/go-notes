package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	pb "play-around/grpc/model"
	"time"
)

type server struct {
	serverId string
	pb.UnimplementedDemoServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("server %v receive message from lb\n", s.serverId)
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func main() {
	grpcOpt1 := grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionAge:      time.Duration(1000) * time.Second,
			MaxConnectionAgeGrace: time.Duration(1000) * time.Second,
			Time:                  time.Duration(1000) * time.Second,
			Timeout:               time.Duration(5) * time.Second,
			MaxConnectionIdle:     time.Duration(100) * time.Second,
		},
	)

	grpcOpt2 := grpc.KeepaliveEnforcementPolicy(
		keepalive.EnforcementPolicy{
			MinTime:             2 * time.Second,
			PermitWithoutStream: true})

	srv := grpc.NewServer(grpcOpt1, grpcOpt2)

	pb.RegisterDemoServiceServer(srv, &server{serverId: "1"})

	fmt.Println("Server is running on port 5577")

	lis, _ := net.Listen("tcp", ":5577")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
