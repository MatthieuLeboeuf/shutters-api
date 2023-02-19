package main

import (
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"log"
	"net/http"
	"time"
)

func up(w http.ResponseWriter, r *http.Request) {
	l, _ := gpiod.RequestLine("gpiochip0", rpi.GPIO3, gpiod.AsOutput(1))
	_ = l.SetValue(0)
	time.Sleep(100 * time.Millisecond)
	_ = l.SetValue(1)
	_ = l.Close()
	w.WriteHeader(http.StatusOK)
}

func down(w http.ResponseWriter, r *http.Request) {
	l, _ := gpiod.RequestLine("gpiochip0", rpi.GPIO2, gpiod.AsOutput(1))
	_ = l.SetValue(0)
	time.Sleep(100 * time.Millisecond)
	_ = l.SetValue(1)
	_ = l.Close()
	w.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/open", up)
	mux.HandleFunc("/close", down)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
