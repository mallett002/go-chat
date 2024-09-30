package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{}
)

func socket(context echo.Context) error {
	ws, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	for {
		// write
		err := ws.WriteMessage(websocket.TextMessage, []byte("howdy, client!"))
		if err != nil {
			context.Logger().Error(err)
		}

		// read
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			context.Logger().Error(err)
		}

		fmt.Printf("messageType: %v\n", msgType)
		fmt.Printf("message: %s\n", msg)
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")

	e.GET("/ws", socket)

	e.Logger.Fatal(e.Start(":1323"))
}
