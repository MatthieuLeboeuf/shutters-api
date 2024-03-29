package main

import (
	"encoding/json"
	"github.com/warthog618/gpiod"
	"net/http"
	"os"
	"time"
)

type Shutter struct {
	Name string `json:"name"`
	Up   int    `json:"up"`
	Down int    `json:"down"`
}

var conf = struct {
	AppToken string    `json:"app_token"`
	AppPort  string    `json:"app_port"`
	Shutters []Shutter `json:"shutters"`
}{}

var shuttersQueue []int

func getShutter(name string) Shutter {
	var shutter Shutter
	for i := 0; i < len(conf.Shutters); i++ {
		if conf.Shutters[i].Name == name {
			shutter = conf.Shutters[i]
		}
	}
	return shutter
}

func pressButton(gpio int) {
	l, _ := gpiod.RequestLine("gpiochip0", gpio, gpiod.AsOutput(1))
	// Press button
	_ = l.SetValue(0)
	time.Sleep(1 * time.Second)
	_ = l.SetValue(1)
	_ = l.Close()
}

func set(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+conf.AppToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// get variables
	_ = r.ParseForm()
	s := r.Form.Get("s")
	p := r.Form.Get("p")

	shutter := getShutter(s)

	// determinate the right button
	gpio := shutter.Down
	if p == "up" {
		gpio = shutter.Up
	}

	// add shutter action to queue
	shuttersQueue = append(shuttersQueue, gpio)

	w.WriteHeader(http.StatusOK)
}

func processQueue() {
	for true {
		for i := 0; i < len(shuttersQueue); i++ {
			pressButton(shuttersQueue[i])
			time.Sleep(1 * time.Second)
		}
		shuttersQueue = []int{}
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	file, _ := os.ReadFile("./config.json")
	err := json.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}

	// start process queue
	go processQueue()

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", set)

	println("Webserver started on 0.0.0.0:" + conf.AppPort)
	err = http.ListenAndServe(":"+conf.AppPort, mux)
	if err != nil {
		panic(err)
	}
}
