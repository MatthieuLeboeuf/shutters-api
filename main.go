package main

import (
	"fmt"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

var percentage int

func set(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	p, _ := strconv.Atoi(r.Form.Get("p"))
	if percentage == p {
		w.WriteHeader(http.StatusOK)
		return
	}
	gpio := rpi.GPIO3 // up
	if p > percentage {
		gpio = rpi.GPIO2 // down
	}
	l, _ := gpiod.RequestLine("gpiochip0", gpio, gpiod.AsOutput(1))
	// Press button
	_ = l.SetValue(0)
	time.Sleep(100 * time.Millisecond)
	_ = l.SetValue(1)
	if p != 0 && p != 100 {
		time.Sleep(time.Duration(17.0/100.0*math.Abs(float64(p-percentage))) * time.Second)
		// Press button
		_ = l.SetValue(0)
		time.Sleep(100 * time.Millisecond)
		_ = l.SetValue(1)
	}
	_ = l.Close()
	percentage = p
	w.WriteHeader(http.StatusOK)
}

func main() {
	percentage = 0
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", set)

	port := "8080"
	fmt.Println("Webserver started on 0.0.0.0:" + port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
