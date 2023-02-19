package main

import (
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

func main() {
	l, err := gpiod.RequestLine("gpiochip0", rpi.GPIO2, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}
	err = l.SetValue(1)
	if err != nil {
		panic(err)
	}
}
