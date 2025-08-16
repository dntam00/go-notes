package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"log"
	"play-around/grpc-ka/pf"
	pb "play-around/grpc/model"
	"play-around/utils"
	"time"
)

func main() {

	ka := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Duration(10) * time.Second,
		Timeout:             time.Duration(3) * time.Second,
		PermitWithoutStream: true,
	})

	dialOpt := grpc.WithConnectParams(grpc.ConnectParams{
		MinConnectTimeout: 500 * time.Millisecond,
	})

	// "name": [{"service": "dnt.DemoService","method": "SayHello"}],
	var retryPolicy = fmt.Sprintf(`{
        "methodConfig": [{
            "name": [],
            "waitForReady": true,
            "retryPolicy": {
                "MaxAttempts": %v,
                "InitialBackoff": "%v",
                "MaxBackoff": "%v",
                "BackoffMultiplier": %v,
                "RetryableStatusCodes": [ "UNAVAILABLE", "UNKNOWN","ABORTED","RESOURCE_EXHAUSTED"]
            }
        }]
    }`, "5", "1s", "2000s", "0.5")

	retry := grpc.WithDefaultServiceConfig(retryPolicy)

	conn, err := grpc.NewClient("127.0.0.1:5577", grpc.WithTransportCredentials(insecure.NewCredentials()), ka, dialOpt, retry)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	testSingleCall(conn)

	utils.Wait()
}

func testSingleCall(conn *grpc.ClientConn) {
	client := pb.NewDemoServiceClient(conn)
	start := time.Now()

	go func() {
		// wait after connection lifetime
		time.Sleep(12 * time.Second)
		timeout, cancelFunc := context.WithTimeout(context.Background(), 120*time.Second)
		_, err := client.SayHello(timeout, &pb.HelloRequest{Name: "world"})
		if err != nil {
			log.Printf("second: could not send after: %v, error: %v\n", time.Since(start), err)
			s, ok := status.FromError(err)
			if ok {
				log.Printf("second: status code: %v, message: %v\n", s.Code(), s.Message())
			} else {
				log.Printf("second: not gRPC status error: %v\n", err)
			}
		} else {
			log.Printf("second: send after: %v\n", time.Since(start))
		}
		cancelFunc()
	}()

	for {
		timeout, cancelFunc := context.WithTimeout(context.Background(), 120*time.Second)
		_, err := client.SayHello(timeout, &pb.HelloRequest{Name: "world"})

		if err != nil {
			log.Printf("could not send after: %v, error: %v\n", time.Since(start), err)
			cancelFunc()
			s, ok := status.FromError(err)
			if ok {
				log.Printf("status code: %v, message: %v\n", s.Code(), s.Message())
			} else {
				log.Printf("not gRPC status error: %v\n", err)
			}
		} else {
			log.Printf("send after: %v\n", time.Since(start))
		}
		time.Sleep(11 * time.Second)
	}
}

func test(c pb.DemoServiceClient) bool {
	pauseCnt := 0

	//ctx, cancelFunc := context.WithTimeout(context.Background(), 1000*time.Second)
	//defer cancelFunc()

	//go func() {
	//	time.Sleep(10 * time.Second)
	//	cancelFunc()
	//}()

	//_, err = c.SayHello(ctx, &pb.HelloRequest{Name: "Hello"})
	//if err != nil {
	//	log.Printf("error could not greet: %v\n", err)
	//}

	//time.Sleep(30 * time.Second)

	//return

	defer func() {
		pf.ApplyRule(pf.DropBlockRule)
	}()

	for {
		if pauseCnt == 2 {
			for i := 0; i < 11; i++ {
				log.Printf("start sleep %v", i)
				time.Sleep(1 * time.Second)
				if i == 7 {
					pf.ApplyRule(pf.BlockRule)
				}
			}
			pauseCnt = 0
		}
		start := time.Now()
		timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		res, err := c.SayHello(timeout, &pb.HelloRequest{Name: "world"})
		if err != nil {
			log.Printf("could not send after: %v, error: %v\n", time.Since(start), err)
			cancelFunc()
			s, ok := status.FromError(err)
			if ok {
				log.Printf("status code: %v, message: %v\n", s.Code(), s.Message())
			} else {
				log.Printf("not gRPC status error: %v\n", err)
			}
			time.Sleep(500 * time.Second)
			return true
		}
		log.Println("receive", res.Message)
		time.Sleep(1 * time.Second)
		pauseCnt++
	}

	////time.Sleep(20 * time.Second)
	//
	//_, err = c.SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	return false
}
