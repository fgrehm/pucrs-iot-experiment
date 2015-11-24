package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

type Buttons interface {
	Run()
	Close()
	Clicks() map[string]int
	ClickedOn(button string)
	ClickSubscribe(listener ClickListener)
}

type ClickListener interface {
	OnClick(button string)
}

type buttons struct {
	clicks    map[string]int
	listeners []ClickListener
	gpioBtns  []*gpioButton
	lastClick time.Time
}

type gpioButton struct {
	number string
	pin    embd.DigitalPin
}

func newButtons() Buttons {
	btns := &buttons{
		clicks:   map[string]int{"0": 0, "1": 0, "2": 0},
		gpioBtns: []*gpioButton{},
	}

	for i := 0; i < 3; i++ {
		btnNum := i
		pinNumber := 22 + i

		log.Println("Set up", btnNum+1, "at pin", pinNumber)

		pin, err := embd.NewDigitalPin(pinNumber)
		if err != nil {
			fmt.Printf("Error in NewDigitalPin (%d)!", pinNumber)
			panic(err)
		}

		if err := pin.SetDirection(embd.In); err != nil {
			fmt.Printf("Error in SetDirection (%d)!", pinNumber)
			panic(err)
		}

		gpioBtn := &gpioButton{number: fmt.Sprintf("%d", btnNum), pin: pin}
		btns.gpioBtns = append(btns.gpioBtns, gpioBtn)
	}

	return btns
}

func (b *buttons) Run() {
	b.ClickSubscribe(b)

	for _, gpioBtn := range b.gpioBtns {
		func (gpioBtn *gpioButton) {
			gpioBtn.pin.Watch(embd.EdgeBoth, func(pin embd.DigitalPin) {
				btnVal, _ := pin.Read()
				if btnVal != embd.Low {
					return
				}

				// HACK: For whatever reason this was triggering a whole lot of clicks
				//       on a button
				if time.Now().Sub(b.lastClick) < (1500 * time.Millisecond) {
					log.Println("Ignoring click on `%s`", gpioBtn.number)
					return
				}

				b.lastClick = time.Now()
				b.ClickedOn(gpioBtn.number)
			})
		}(gpioBtn)
	}
}

func (b *buttons) Close() {
	log.Println("Closing buttons...")
	for _, gpioBtn := range b.gpioBtns {
		gpioBtn.pin.Close()
	}

	log.Println("Closing GPIO...")
	embd.CloseGPIO()
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
