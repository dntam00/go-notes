package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(serveWebSocket)
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "PUT", "GET", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%+v", 5001),
		Handler: handlers.CORS(credentials, methods, origins)(r),
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	_ = conn.Close()
}
