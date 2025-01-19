package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	pb "play-around/grpc/model"
	"strings"
	"sync"
)

const (
	address = "static:///localhost:50051,localhost:50052,localhost:50053"
)

func init() {
	resolver.Register(&StaticBuilder{})
}

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`)}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		panic(err)
	}

	c := pb.NewDemoServiceClient(conn)
	for i := 0; i < 10; i++ {
		_, err = c.SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	}
}

type StaticBuilder struct{}

func (sb *StaticBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	endpoints := strings.Split(target.Endpoint(), ",")

	r := &StaticResolver{
		endpoints: endpoints,
		cc:        cc,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (sb *StaticBuilder) Scheme() string {
	return "static"
}

type StaticResolver struct {
	endpoints []string
	cc        resolver.ClientConn
	sync.Mutex
}

func (r *StaticResolver) ResolveNow(opts resolver.ResolveNowOptions) {
	r.Lock()
	r.doResolve()
	r.Unlock()
}

func (r *StaticResolver) Close() {
}

func (r *StaticResolver) doResolve() {
	var addrs []resolver.Address
	for i, addr := range r.endpoints {
		addrs = append(addrs, resolver.Address{
			Addr:       addr,
			ServerName: fmt.Sprintf("instance-%d", i+1),
		})
	}

	newState := resolver.State{
		Addresses: addrs,
	}

	err := r.cc.UpdateState(newState)
	if err != nil {
		log.Printf("failed to update state: %v", err)
	}
}
