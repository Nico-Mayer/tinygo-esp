package main

import (
	"blinky-esp/internal/display"
	_ "embed"
	"fmt"
	"io"
	"machine"
	"net/http"
	"sync/atomic"
	"time"

	"tinygo.org/x/drivers/netdev"
	nl "tinygo.org/x/drivers/netlink"
	link "tinygo.org/x/espradio/netlink"
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

	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	serial := machine.Serial
	serial.Configure(machine.UARTConfig{BaudRate: 115200})
	serial.Write([]byte("HTTP Server\r\n"))

	display := display.NewDisplay()

	radioLink := link.Esplink{
		ArenaPoolSize: 48 * 1024,
	}
	netdev.UseNetdev(&radioLink)

	serial.Write([]byte("Connecting to Wi-Fi...\r\n"))
	display.Show("Connecting...")
	err := radioLink.NetConnect(&nl.ConnectParams{
		Ssid:       ssid,
		Passphrase: pw,
	})
	if err != nil {
		serial.Write([]byte("Failed to connect to wifi\r\n"))
		return
	}

	serial.Write([]byte("Connected!\r\n"))
	display.Show("Connected!")

	time.Sleep(5 * time.Second)

	addr, _ := radioLink.Addr()
	host := addr.String()
	serial.Write([]byte("Server: http://"))
	serial.Write([]byte(host))
	serial.Write([]byte(":8080\r\n"))

	http.Handle("/", logRequest(root))
	http.Handle("/led/on", logRequest(ledOn))
	http.Handle("/led/off", logRequest(ledOff))
	http.Handle("/status", logRequest(status))

	serial.Write([]byte("Starting server...\r\n"))

	go func() {
		for {
			display.Show(fmt.Sprintf("LED: %v", led.Get()))
			time.Sleep(time.Second)
		}
	}()

	err = http.ListenAndServe(host+":8080", nil)
	if err != nil {
		serial.Write([]byte("Server error: "))
		serial.Write([]byte(err.Error()))
		serial.Write([]byte("\r\n"))
	}
}

func logRequest(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serial := machine.Serial
		serial.Write([]byte(r.Method))
		serial.Write([]byte(" "))
		serial.Write([]byte(r.URL.Path))
		serial.Write([]byte("\r\n"))
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
