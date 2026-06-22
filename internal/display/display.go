package display

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

const (
	SDA_PIN = machine.GPIO8
	SCL_PIN = machine.GPIO9
)

type Display struct {
	device *ssd1306.Device
}

func NewDisplay() *Display {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		SDA: SDA_PIN,
		SCL: SCL_PIN,
	})

	time.Sleep(2 * time.Second)

	device := ssd1306.NewI2C(i2c)
	device.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})

	device.Display()
	device.ClearBuffer()
	return &Display{
		device: device,
	}
}

func (d *Display) Show(s string) {
	d.device.ClearBuffer()
	tinyfont.WriteLine(d.device, &freemono.Regular9pt7b,
		0, 20, s, color.RGBA{255, 255, 255, 255})
	d.device.Display()
}
