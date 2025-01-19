package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"play-around/dict/html"
	"play-around/dict/model"
	"play-around/dict/monitor"
	sql "play-around/dict/store"
	worker2 "play-around/dict/worker"
	"runtime"
	"syscall"
	"time"
)

type entry struct {
	word    string
	content string
}

const (
	delimiter = "</>\r\n"
	batchSize = 1000
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	monitor.Start()
	go monitor.ServeMetrics()

	start := time.Now()

	numWorkers := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(numWorkers)

	// Open the file
	file, err := os.Open("/Users/dntam/Projects/golang/play-around/dict/ldoce.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Create a new scanner

	reader := bufio.NewReader(file)

	store := make([]string, 0)
	messages := make([]model.Entry, batchSize)

	connection, err := sql.InitDbConnection()

	if err != nil {
		log.Fatalf("failed to connect to db: %s", err)
	}

	processor := html.New(connection)

	w := worker2.New(numWorkers, 10000, connection, processor)

	//w.StartWorker()
	w.StartWorkerEntries()

	index := 0

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading file", err)
			break
		}
		if text == delimiter {
			en := model.Entry{
				Word:    store[0],
				Content: store[1],
			}
			//w.AddEntry(en)
			store = make([]string, 0)
			messages[index] = en
			index++
			if index == batchSize {
				w.AddEntries(messages)
				index = 0
				messages = make([]model.Entry, batchSize)
			}
			continue
		}
		store = append(store, text)
	}

	fmt.Println("finish reading dictionary, take:", time.Since(start).Seconds())
	<-ctx.Done()
}
