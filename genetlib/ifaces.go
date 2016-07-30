// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"net"
)

type GenetIface struct {
	Index        int
	Mtu          int
	Name         string
	Hardwareaddr string
	Flags        string
	Addr         string
}

func GetIfaces() ([]GenetIface, error) {

	nifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	genets := make([]GenetIface, 0)

	for z, nif := range nifs {
		addrs := ""
		paddrs, _ := nifs[z].Addrs()
		for _, addr := range paddrs {
			if addrs != "" {
				addrs += "|"
			}
			addrs = addrs + addr.String()
		}
		mygenetiface := GenetIface{nif.Index, nif.MTU, nif.Name, nif.HardwareAddr.String(), nif.Flags.String(), addrs}
		genets = append(genets, mygenetiface)
	}

	return genets, nil
}
