package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/tarm/serial"
)

type Leds interface {
	OnClick(button string)
}

type leds struct {
	btns Buttons
	hub  *Hub
}

func newLeds(btns Buttons, hub *Hub) Leds {
	leds := &leds{
		btns: btns,
		hub:  hub,
	}
	btns.ClickSubscribe(leds)
	return leds
}

func (l *leds) OnClick(button string) {
	totalClicks := l.btns.Clicks()[button]

	l.hub.Broadcast <- fmt.Sprintf(`{"event":"processing","data":{"button":%s}}`, button)

	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	btnNum, err := strconv.Atoi(button)
	n, err := s.Write([]byte{byte(btnNum), byte(totalClicks)})
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response from arduino: %q", buf[:n])

	l.hub.Broadcast <- fmt.Sprintf(`{"event":"click","data":{"button":%s,"count":%d}}`, button, totalClicks)
}
