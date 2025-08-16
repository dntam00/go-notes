package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func doRequest(client *http.Client, req *http.Request) {
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error call api: ", err)
	}

	defer func() {
		if res != nil {
			_, _ = io.Copy(io.Discard, res.Body)
			_ = res.Body.Close()
		}
	}()
}

func main() {
	fmt.Println("start test")
	client := NewHttpClient()

	url := "http://127.0.0.1:7888/endpoint"

	for i := 0; i < 5; i++ {
		//go func() {
		timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)

		req, err := http.NewRequestWithContext(timeout, http.MethodPost, url, nil)
		if err != nil {
			fmt.Println("error: ", err)
		}
		req.Header.Set("Content-Type", "application/json")

		doRequest(client, req)
		cancelFunc()
		//}()
	}

	//go func() {
	//	time.Sleep(9 * time.Second)
	//	fmt.Println("start second request")
	//	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, nil)
	//	if err != nil {
	//		fmt.Println("error: ", err)
	//	}
	//	req.Header.Set("Content-Type", "application/json")
	//
	//	doRequest(client, req)
	//}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}

func NewHttpClient() *http.Client {
	transport := NewTransport()

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(200) * time.Second,
	}
}

func NewTransport() http.RoundTripper {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxConnsPerHost = 1
	t.IdleConnTimeout = 30 * time.Second
	t.MaxIdleConns = 1
	t.MaxIdleConnsPerHost = 1
	return t
}
