package main

import (
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"time"
)

func main() {
	c, _ := gpiod.NewChip("gpiochip0", gpiod.WithConsumer("myapp"))
	l, _ := c.RequestLine(rpi.GPIO2)
	_ = l.SetValue(0)
	time.Sleep(100 * time.Millisecond)
	_ = l.SetValue(1)
	_ = l.Close()
}
