package main

import (
	"encoding/json"
	"github.com/warthog618/gpiod"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Shutter struct {
	Name  string `json:"name"`
	Up    int    `json:"up"`
	Down  int    `json:"down"`
	Total int    `json:"total"`
}

type Config struct {
	Api struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	} `json:"api"`
	Shutters []Shutter `json:"shutters"`
}

type apiEntity struct {
	Attributes struct {
		Percentage int `json:"current_position"`
	} `json:"attributes"`
}

var conf = Config{}

func getShutter(name string) Shutter {
	var shutter Shutter
	for i := 0; i < len(conf.Shutters); i++ {
		if conf.Shutters[i].Name == name {
			shutter = conf.Shutters[i]
		}
	}
	return shutter
}

func getPercentage(shutter Shutter) int {
	req, _ := http.NewRequest(
		"GET",
		conf.Api.Url+"/states/cover.shutter_"+shutter.Name,
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+conf.Api.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	entity := apiEntity{}
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &entity)
	defer resp.Body.Close()

	return entity.Attributes.Percentage
}

func pressButton(gpio int) {
	l, _ := gpiod.RequestLine("gpiochip0", gpio, gpiod.AsOutput(1))
	// Press button
	_ = l.SetValue(0)
	time.Sleep(100 * time.Millisecond)
	_ = l.SetValue(1)
	_ = l.Close()
}

func set(w http.ResponseWriter, r *http.Request) {
	// get variables
	_ = r.ParseForm()
	s := r.Form.Get("s")
	p, _ := strconv.Atoi(r.Form.Get("p"))

	shutter := getShutter(s)

	// fetch current percentage from api
	actual := getPercentage(shutter)

	// determinate the right button
	gpio := shutter.Down
	if p > actual {
		gpio = shutter.Up
	}

	// start shutter
	pressButton(gpio)

	// wait for the action
	if p != 0 && p != 100 {
		time.Sleep(time.Duration(float64(shutter.Total)/100.0*math.Abs(float64(p-actual))) * time.Second)
		pressButton(gpio)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	configFile, _ := os.ReadFile("./config.json")
	err := json.Unmarshal(configFile, &conf)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", set)

	port := "8080"
	println("Webserver started on 0.0.0.0:" + port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
