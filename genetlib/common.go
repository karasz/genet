// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"bufio"
	"log"
	"math"
	"net"
	"os"
	"strings"
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

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}

func MakeStdDev(series []float64, avg float64) float64 {
	var sumsquares float64

	if len(series) <= 1 {
		return 0
	}

	for _, s := range series {
		sumsquares += math.Pow(s-avg, 2)
	}
	varance := sumsquares / float64(len(series)-1)
	return math.Sqrt(varance)
}
