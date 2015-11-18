package main

import (
	"os"
	"fmt"
	"log"
	"time"
)

type Buttons interface {
	Run()
	Clicks() map[string]int
	ClickedOn(button string)
	ClickSubscribe(listener ClickListener)
}

type ClickListener interface {
	OnClick(button string)
}

type buttons struct {
	gpioRoot  string
	hub       *Hub
	clicks    map[string]int
	listeners []ClickListener
}

func newButtons(gpioRoot string, hub *Hub) Buttons {
	btns := &buttons{
		gpioRoot: gpioRoot,
		hub:      hub,
		clicks:   map[string]int{"0": 0, "1": 0, "2": 0},
	}
	btns.ClickSubscribe(btns)
	return btns
}

func (b *buttons) Run() {
	for {
		for i := 0; i < 3; i++ {
			path := fmt.Sprintf("/sys/class/gpio/gpio%d/value", 25 - i)
			log.Println("Reading", i, "at", path)
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			value := make([]byte, 1)
			if _, err = file.Read(value); err != nil {
				log.Fatal(err)
			}
			file.Close()
			log.Println("Found", string(value), "at", path)
			if string(value) == "0" {
				b.ClickedOn(fmt.Sprintf("%d", i))
			}
		}
		time.Sleep(1300 * time.Second)
	}
}

func (b *buttons) Clicks() map[string]int {
	return b.clicks
}

func (b *buttons) ClickSubscribe(listener ClickListener) {
	b.listeners = append(b.listeners, listener)
}

func (b *buttons) ClickedOn(button string) {
	for _, l := range b.listeners {
		l.OnClick(button)
	}
}

func (b *buttons) OnClick(button string) {
	log.Println("Clicked on", button)
	b.clicks[button] += 1
}
