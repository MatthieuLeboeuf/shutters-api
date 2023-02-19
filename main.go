package main

import (
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"time"
)

func main() {
	l, err := gpiod.RequestLine("gpiochip0", rpi.GPIO2)
	if err != nil {
		panic(err)
	}
	err = l.SetValue(0)
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)
	err = l.SetValue(1)
	if err != nil {
		panic(err)
	}
	err = l.Close()
	if err != nil {
		panic(err)
	}
}
