package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

var (
	hub = &Hub{
		Broadcast:   make(chan string),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
		Connections: make(map[*Connection]bool),
	}
	btns Buttons
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

	// Detect clicks on buttons
	btns = newButtons("gpio", hub)
	go btns.Run()

	// Control the leds based on clicks
	newLeds(btns, hub)

	// Start server
	fmt.Println("Starting server on port " + port)
	e.Run(":" + port)
}

func buttonsStateHandler(c *echo.Context) error {
	response := btns.Clicks()
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
	go btns.ClickedOn(c.Param("button"))
	return c.NoContent(http.StatusCreated)
}
