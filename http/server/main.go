package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Event struct {
	AppID       string    `json:"app_id"`
	UserIDByApp string    `json:"user_id_by_app"`
	EventName   string    `json:"event_name"`
	Timestamp   string    `json:"timestamp"`
	Sender      Sender    `json:"sender"`
	Recipient   Recipient `json:"recipient"`
	Message     Message   `json:"message"`
	Source      string    `json:"source"`
	Follower    Follower  `json:"follower"`
}

type Follower struct {
	ID string `json:"id"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	MsgID string `json:"msg_id"`
	Text  string `json:"text"`
}

func main() {
	headersOk := handlers.AllowedHeaders([]string{
		"sec-ch-ua", "x-owner-content", "sec-ch-ua-mobile", "User-Agent",
		"Accept", "Referer", "device", "sec-ch-ua-platform",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r := mux.NewRouter()
	r.HandleFunc("/endpoint", func(writer http.ResponseWriter, request *http.Request) {

		time.Sleep(3 * time.Second)
		fmt.Println("request header", request.Header)

		var event Event

		err := json.NewDecoder(request.Body).Decode(&event)
		if err != nil {
			http.Error(writer, "Invalid request body", http.StatusBadRequest)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// encode event to json
		eventJson, err := json.Marshal(event)

		fmt.Printf("body: %+v", string(eventJson))

		//hostname, err := os.Hostname()
		//if err != nil {
		//	fmt.Println("error getting hostname: ", err)
		//}
		writer.WriteHeader(http.StatusOK)
		//_, err = writer.Write([]byte(hostname))
		//if err != nil {
		//	fmt.Println("write error: ", err)
		//}
	}).Methods("POST")

	srv := &http.Server{
		Addr:        ":" + "7888",
		Handler:     handlers.CORS(originsOk, headersOk, methodsOk)(r),
		IdleTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("failed to start server")
		return
	}
}
