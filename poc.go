package main

import (
	"fmt"
	"flag"
	"log"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
    fmt.Println("Hello, World!")
	flag.Parse()
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// http request handler: Think I'll need to use a server like echo/chi
func handler(writer http.ResponseWriter, req *http.Request) {
    conn, err := upgrader.Upgrade(writer, req, nil)

    if err != nil {
        log.Println(err)
        return
    }

    // Use conn to send and receive messages:
	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}

}
