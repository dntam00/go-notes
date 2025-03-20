package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func main() {
	headersOk := handlers.AllowedHeaders([]string{
		"sec-ch-ua", "x-owner-content", "sec-ch-ua-mobile", "User-Agent",
		"Accept", "Referer", "device", "sec-ch-ua-platform",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r := mux.NewRouter()
	r.HandleFunc("/endpoint", func(writer http.ResponseWriter, request *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("error getting hostname: ", err)
		}
		_, err = writer.Write([]byte(hostname))
		if err != nil {
			fmt.Println("write error: ", err)
		}
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
