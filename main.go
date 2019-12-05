package main

import (
	"log"
	"net/http"
	"web/websocket/actions"
	"web/websocket/client"

	"golang.org/x/net/websocket"
)

func main() {
	server := client.NewServer(actions.NewHandlers(actions.Routes))
	srv := websocket.Server{
		Handler: server.Handle,
	}
	log.Println("websocket server listens on port 1234")
	log.Fatalln(http.ListenAndServe(":12345", srv))

}
