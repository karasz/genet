// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
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

func CompareVersions(v1 string, v2 string) (int64, error) {
	if strings.Index(v1, ".") == -1 || strings.Index(v2, ".") == -1 {
		return 0, fmt.Errorf("No dot in version strings %s, %s", v1, v2)
	}

	vs1 := strings.Split(v1, ".")
	vs2 := strings.Split(v2, ".")

	if _, err := strconv.Atoi(vs1[0]); err != nil {
		return 0, fmt.Errorf("Version %s does not begin with a number", v1)
	}
	if _, err := strconv.Atoi(vs2[0]); err != nil {
		return 0, fmt.Errorf("Version %s does not begin with a number", v2)
	}

	var slic []string

	if len(vs1) >= len(vs2) {
		slic = vs2
	} else {
		slic = vs1
	}

	for i, _ := range slic {
		if vs1[i] > vs2[i] {
			return 1, nil
			break
		}
		if vs1[i] < vs2[i] {
			return -1, nil
			break
		}
		if vs1[i] == vs2[i] {
			continue
		}
	}
	return 0, fmt.Errorf("No idea what happened :/")

}

func LookupLinuxUserById(id int) (string, error) {
	usermap := make(map[int]string)

	lines, err := ReadLines("/etc/passwd")
	if err != nil {
		return "", err
	}

	for _, line := range lines {
		data := strings.Split(line, ":")
		uid, err := strconv.Atoi(data[2])

		if err != nil {
			return "", err
		}

		usermap[uid] = data[0]
	}

	if usermap[id] != "" {
		return usermap[id], nil
	}
	return "", fmt.Errorf("No User with UID = %d", id)
}
