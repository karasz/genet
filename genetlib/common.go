// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"bufio"
	"net"
	"os"
)

func ReadLines(filename string) ([]string, error) {

	file, err := os.Open(filename)

	if err != nil {
		return []string{""}, err
	}

	defer file.Close()

	var ret []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	return ret, scanner.Err()
}

func GetIpFromName(ifaceName string) ([]string, error) {
	var result []string
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return result, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return result, err
	}
	for _, addr := range addrs {
		ip, _, _ := net.ParseCIDR(addr.String())
		result = append(result, ip.String())
	}
	return result, nil

}
