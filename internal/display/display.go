package display

import (
	"image/color"
	"machine"
	"strings"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

type Display struct {
	device *ssd1306.Device
}

func NewDisplay(sda_pin, scl_pin machine.Pin) *Display {
	i2c := machine.I2C0
	i2c.Configure(machine.I2CConfig{
		SDA: sda_pin,
		SCL: scl_pin,
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

const lineHeight = 16 // freemono 9pt ~13px tall + spacing

func (d *Display) Show(args ...string) {
	var builder strings.Builder

	for _, s := range args {
		builder.WriteString(s)
	}

	d.device.ClearBuffer()
	white := color.RGBA{255, 255, 255, 255}
	for i, line := range strings.Split(builder.String(), "\n") {
		y := int16(lineHeight * (i + 1))
		tinyfont.WriteLine(d.device, &freemono.Regular9pt7b, 0, y, line, white)
	}
	d.device.Display()
}

func (d *Display) ShowAndLog(args ...string) {
	var builder strings.Builder
	for _, s := range args {
		builder.WriteString(s)
	}

	println(builder.String())
	d.Show(args...)
}
