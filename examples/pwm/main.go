package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

const (
	pollRate = time.Millisecond * 500
)

func main() {
	display := initDisplay()

	machine.InitADC()
	pot := machine.ADC{Pin: machine.ADC4}
	pot.Configure(machine.ADCConfig{})

	pwm := machine.PWM0
	pinA := machine.GPIO17
	// Configure the PWM with the given period.
	err := pwm.Configure(machine.PWMConfig{
		Period: 16384e3, // 16.384ms
	})
	if err != nil {
		println("failed to configure PWM")
		return
	}

	ch, err := pwm.Channel(pinA)
	if err != nil {
		println("failed to get PWM channel")
		return
	}

	var prev uint16

	for {
		val := pot.Get()
		msg := fmt.Sprintf("V: %d\n", val)

		machine.Serial.Write([]byte(msg))

		top := pwm.Top()
		duty := uint32(val) * top / 65535
		pwm.Set(ch, duty)

		if val != prev {
			display.ClearBuffer()
			tinyfont.WriteLine(display, &freemono.Oblique9pt7b,
				0, 20, msg, color.RGBA{255, 255, 255, 255})
			display.Display()
		}

		prev = val
		time.Sleep(pollRate)
	}
}

func initDisplay() *ssd1306.Device {
	scl := machine.GPIO9
	sda := machine.GPIO8

	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		SDA: sda,
		SCL: scl,
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
