package main

import (
	"blinky-esp/internal/display"
	"blinky-esp/internal/wifi"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"machine"
	"net/http"
	"time"
)

const (
	ssid = ""
	pw   = ""
)

//go:embed index.html
var html string

type RGBLed struct {
	pwm      *machine.LEDCPWM
	channels [3]uint8
	rgb      [3]uint8
}

func NewRGBLed(rPin, gPin, bPin machine.Pin) (*RGBLed, error) {
	l := &RGBLed{pwm: machine.PWM0}
	if err := l.pwm.Configure(machine.PWMConfig{Period: 1e9 / 5000}); err != nil {
		return nil, err
	}
	for i, pin := range [3]machine.Pin{rPin, gPin, bPin} {
		ch, err := l.pwm.Channel(pin)
		if err != nil {
			return nil, err
		}
		l.channels[i] = ch
	}
	return l, nil
}

func (l *RGBLed) SetColor(r, g, b uint8) {
	l.rgb = [3]uint8{r, g, b}
	top := l.pwm.Top()
	for i, v := range l.rgb {
		duty := uint32(v) * top / 255
		l.pwm.Set(l.channels[i], duty)
	}
}

func (l *RGBLed) Off() { l.SetColor(0, 0, 0) }
func (l *RGBLed) State() string {
	return fmt.Sprintf("R: %d\nG: %d\nB: %d", l.rgb[0], l.rgb[1], l.rgb[2])
}

func main() {
	display := display.NewDisplay()

	led, err := NewRGBLed(machine.GPIO7, machine.GPIO6, machine.GPIO5)
	if err != nil {
		println("led init:", err.Error())
		return
	}
	led.SetColor(255, 0, 250)

	display.ShowAndLog("Connecting to wifi...")
	radioLink, err := wifi.Connect(ssid, pw)
	if err != nil {
		println("Failed to connect wifi.")
		return
	}

	addr, _ := radioLink.Addr()
	host := addr.String()
	display.ShowAndLog("Server: http://", host, ":8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, html)
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		var c struct{ R, G, B uint8 }
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		led.SetColor(c.R, c.G, c.B)
		w.WriteHeader(http.StatusNoContent)
	})

	go func() {
		for {
			display.Show(led.State())
			time.Sleep(500 * time.Millisecond)
		}
	}()

	if err := http.ListenAndServe(host+":8080", nil); err != nil {
		println("Server error:", err.Error())
	}
}
