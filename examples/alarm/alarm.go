package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
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
	display := initDisplay()

	var lastMotion time.Time
	alarmOn := false

	for {
		if motionSensor.Get() {
			lastMotion = time.Now()
			if !alarmOn {
				buzzer.High()
				showAlarm(display)
				alarmOn = true
				serial.Write([]byte("Motion detected: buzzer on\n"))
			}
		} else if alarmOn && time.Since(lastMotion) >= holdTime {
			buzzer.Low()
			alarmOn = false
			serial.Write([]byte("Buzzer off: hold time elapsed\n"))
			clearScreen(display)
		}

		time.Sleep(pollInterval)
	}
}

func showAlarm(display *ssd1306.Device) {
	display.ClearBuffer()

	tinyfont.WriteLine(display, &freemono.Regular9pt7b,
		0, 20, "Alarm!", color.RGBA{255, 255, 0, 255})

	display.Display()
}

func clearScreen(display *ssd1306.Device) {
	display.ClearBuffer()
	display.Display()
}

func initDisplay() *ssd1306.Device {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		SDA: machine.GPIO8,
		SCL: machine.GPIO9,
	})

	time.Sleep(2 * time.Second)

	display := ssd1306.NewI2C(i2c)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})

	display.Display()
	display.ClearBuffer()

	return display
}
