package wifi

import (
	"tinygo.org/x/drivers/netdev"
	nl "tinygo.org/x/drivers/netlink"
	link "tinygo.org/x/espradio/netlink"
)

func Connect(ssid, pw string) (*link.Esplink, error) {
	radioLink := link.Esplink{
		ArenaPoolSize: 48 * 1024,
	}
	netdev.UseNetdev(&radioLink)

	err := radioLink.NetConnect(&nl.ConnectParams{
		Ssid:       ssid,
		Passphrase: pw,
	})

	if err != nil {
		return nil, err
	}

	return &radioLink, nil
}
