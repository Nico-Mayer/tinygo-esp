package main

import (
	"blinky-esp/internal/display"
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/dht"
)

func logf(format string, args ...any) {
	fmt.Fprintf(machine.Serial, format+"\n", args...)
}

// This example only works with the normal esp32 and not esp32s3 because the humid sensor is not working on the s3 version

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})

	disp := display.NewDisplay(machine.GPIO21, machine.GPIO2) // check befor flashing if those are the scl and sda pins
	sensor := dht.New(machine.GPIO4, dht.DHT11)

	disp.ShowAndLog("Connecting...")

	for {
		if err := sensor.ReadMeasurements(); err != nil {
			logf("read error: %s", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}

		temp, _ := sensor.TemperatureFloat(dht.C)
		hum, _ := sensor.HumidityFloat()

		logf("Temp: %.1f C  Hum: %.0f %%", temp, hum)
		disp.Show(fmt.Sprintf("Temp: %.0f C\nHum: %.0f %%", temp, hum))

		time.Sleep(2 * time.Second)
	}
}
