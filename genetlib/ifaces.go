package genetlib

import (
	"net"
)

type genetiface struct {
	index        int
	mtu          int
	name         string
	hardwareaddr string
	flags        string
	addr         string
}

func GetIfaces() ([]genetiface, error) {

	nifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	genets := make([]genetiface, 0)

	for z, nif := range nifs {
		addrs := ""
		paddrs, _ := nifs[z].Addrs()
		for _, addr := range paddrs {
			if addrs != "" {
				addrs += "|"
			}
			addrs = addrs + addr.String()
		}
		mygenetiface := genetiface{nif.Index, nif.MTU, nif.Name, nif.HardwareAddr.String(), nif.Flags.String(), addrs}
		genets = append(genets, mygenetiface)
	}

	return genets, nil
}
