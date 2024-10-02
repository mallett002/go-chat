package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var connectedClients = make(map[*websocket.Conn]bool)

func socket(context echo.Context) error {
	ws, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()


	// For broadcasting to all clients, add this client to the connected clients list
	connectedClients[ws] = true

	for {
		// read msg from user
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			context.Logger().Error(err)
			delete(connectedClients, ws) // remove this client from the list
			break
		}
		
		fmt.Printf("User sent: %s\n", msg)

		// // write back to same user
		// err = ws.WriteMessage(msgType, msg)
		// if err != nil {
		// 	context.Logger().Error(err)
		// 	break
		// }

		for client := range connectedClients {
			if err := client.WriteMessage(msgType, msg); err != nil {
				context.Logger().Error(err)
				client.Close() // close connection
				delete(connectedClients, client) // remove client from list
			}
		}
	}

	return nil
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")

	e.GET("/ws", socket)

	e.Logger.Fatal(e.Start(":1323"))
}
