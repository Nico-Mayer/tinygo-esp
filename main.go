package main

import (
	"machine"
	"time"
)

const (
	buzzer       = machine.GPIO26
	motionSensor = machine.GPIO27

	holdTime     = 500 * time.Millisecond
	pollInterval = 50 * time.Millisecond
)

func main() {
	serial := machine.Serial

	buzzer.Configure(machine.PinConfig{Mode: machine.PinOutput})
	motionSensor.Configure(machine.PinConfig{Mode: machine.PinInput})

	var lastMotion time.Time
	alarmOn := false

	for {
		if motionSensor.Get() {
			lastMotion = time.Now()
			if !alarmOn {
				buzzer.High()
				alarmOn = true
				serial.Write([]byte("Motion detected: buzzer on\n"))
			}
		} else if alarmOn && time.Since(lastMotion) >= holdTime {
			buzzer.Low()
			alarmOn = false
			serial.Write([]byte("Buzzer off: hold time elapsed\n"))
		}

		time.Sleep(pollInterval)
	}
}
