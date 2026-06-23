package main

import (
	"blinky-esp/internal/display"
	"blinky-esp/internal/wifi"
	_ "embed"
	"fmt"
	"io"
	"machine"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	ssid = ""
	pw   = ""
	led  = machine.GPIO17
)

//go:embed index.html
var html string
var ledState atomic.Bool

func main() {
	display := display.NewDisplay()
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	display.ShowAndLog("Connecting to wifi...")
	radioLink, err := wifi.Connect(ssid, pw)
	if err != nil {
		println("Faild to connect wifi.")
		return
	}

	display.ShowAndLog("Connected!")

	time.Sleep(5 * time.Second)

	addr, _ := radioLink.Addr()
	host := addr.String()
	display.ShowAndLog("Server: http://", host, ":8080")

	http.Handle("/", logRequest(root))
	http.Handle("/led/on", logRequest(ledOn))
	http.Handle("/led/off", logRequest(ledOff))
	http.Handle("/status", logRequest(status))

	display.ShowAndLog("Starting server...")

	go func() {
		for {
			display.Show(fmt.Sprintf("LED: %v", led.Get()))
			time.Sleep(time.Second)
		}
	}()

	err = http.ListenAndServe(host+":8080", nil)
	if err != nil {
		println("Server error: ", err.Error())
	}
}

func logRequest(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.Method, " ", r.URL.Path)
		h(w, r)
	})
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, html)
}

func ledOn(w http.ResponseWriter, r *http.Request) {
	setLED(true)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "LED ON")
}

func ledOff(w http.ResponseWriter, r *http.Request) {
	setLED(false)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "LED OFF")
}

func status(w http.ResponseWriter, r *http.Request) {
	stateStr := "OFF"
	if ledState.Load() {
		stateStr = "ON"
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"led":"`+stateStr+`","uptime":"`+getUptime()+`"}`)
}

func getUptime() string {
	return "unknown" // Simplified for example
}

func setLED(on bool) {
	if on {
		led.High()
	} else {
		led.Low()
	}
	ledState.Store(on)
}
