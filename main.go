package main

import (
	"machine"
	"time"
)

const (
	DELAY = 200
)

func main() {
	// Initialize serial for output
	serial := machine.Serial
	serial.Configure(machine.UARTConfig{BaudRate: 115200})

	// Initialize LED
	led := machine.GPIO2
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	serial.Write([]byte("LED Blink Example\r\n"))

	for {
		serial.Write([]byte("LED ON\r\n"))
		led.High()
		time.Sleep(time.Millisecond * DELAY)

		serial.Write([]byte("LED OFF\r\n"))
		led.Low()
		time.Sleep(time.Millisecond * DELAY)
	}
}
