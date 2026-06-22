package main

import (
	"machine"
)

const (
	DELAY = 1000
	btn   = machine.GPIO4
	led   = machine.GPIO5
)

func main() {
	// Configure LED pin for your board:
	// ESP32/ESP32-S3: GPIO2
	// ESP32-C3: GPIO2
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	btn.Configure(machine.PinConfig{Mode: machine.PinInput})

	for {
		if btn.Get() {
			led.High()
		} else {
			led.Low()
		}
	}
}
