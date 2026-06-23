package wifi

import (
	"fmt"
	"time"

	"tinygo.org/x/drivers/netdev"
	nl "tinygo.org/x/drivers/netlink"
	link "tinygo.org/x/espradio/netlink"
)

const (
	MAX_RETRY_COUNT = 7
)

func Connect(ssid, pw string) (*link.Esplink, error) {
	radioLink := link.Esplink{
		ArenaPoolSize: 48 * 1024,
	}
	netdev.UseNetdev(&radioLink)

	for i := range MAX_RETRY_COUNT {
		err := radioLink.NetConnect(&nl.ConnectParams{
			Ssid:       ssid,
			Passphrase: pw,
		})

		if err != nil {
			fmt.Printf("Wifi Connection failed (try:%d): retry...", i+1)
		} else if err != nil && i-1 == MAX_RETRY_COUNT {
			return nil, err
		} else if err == nil {
			break
		}

		time.Sleep(time.Second)
	}

	return &radioLink, nil
}
