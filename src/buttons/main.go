package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

var (
	clicks = map[string]int{"0": 0, "1": 0, "2": 0}
	hub    = &Hub{
		Broadcast:   make(chan string),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
		Connections: make(map[*Connection]bool),
	}
)

func main() {
	port := "8080"

	e := echo.New()
	e.SetDebug(true)

	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.StripTrailingSlash()

	assetHandler := http.FileServer(assetFS())

	e.Get("/clicks", buttonsStateHandler)
	e.Post("/click/:button", buttonClicker)
	e.WebSocket("/ws", socketHandler)
	e.Get("/*", func(c *echo.Context) error {
		assetHandler.ServeHTTP(c.Response().Writer(), c.Request())
		return nil
	})

	// Start the "message hub"
	go hub.Run()

	// Start server
	fmt.Println("Starting server on port " + port)
	e.Run(":" + port)
}

func buttonsStateHandler(c *echo.Context) error {
	response := []int{clicks["0"], clicks["1"], clicks["2"]}
	return c.JSON(http.StatusOK, response)
}

func socketHandler(c *echo.Context) error {
	var err error

	ws := c.Socket()
	conn := &Connection{Send: make(chan string, 256), WS: ws, Hub: hub}
	conn.Hub.Register <- conn
	defer func() { hub.Unregister <- conn }()
	conn.Writer()

	return err
}

func buttonClicker(c *echo.Context) error {
	go func() {
		hub.Broadcast <- fmt.Sprintf(`{"event":"processing","data":{"button":%s}}`, c.Param("button"))
		time.Sleep(1500 * time.Millisecond)
		clicks[c.Param("button")] += 1
		hub.Broadcast <- fmt.Sprintf(`{"event":"click","data":{"button":%s,"count":%d}}`, c.Param("button"), clicks[c.Param("button")])
	}()
	return c.NoContent(http.StatusCreated)
}
