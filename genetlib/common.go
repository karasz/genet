// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

func ReadLinesNoFrills(filename string, hdrLines int, comm string) ([]string, error) {

	z, err := ReadLines(filename)
	if err != nil {
		return []string{""}, err
	}
	z = z[hdrLines:]
	if comm != "" {
		for i, line := range z {
			if strings.HasPrefix(line, comm) {
				z = append(z[:i], z[i+1:]...)
			}
		}
	}
	return z, err
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

func GetValFromFile(filename string, scase int) string {
	z, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	switch scase {
	case 0:
		return strings.ToUpper(string(z[:len(z)-1]))
	case 1:
		return strings.ToLower(string(z[:len(z)-1]))
	default:
		return string(z[:len(z)-1])
	}
}

func GetFieldValFromFile(filename string, bkfld int, bkfldval string, fld int) string {
	lines, err := ReadLinesNoFrills(filename, 0, "#")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var result string
	for _, line := range lines {
		if strings.Fields(line)[bkfld] == bkfldval {
			result = strings.Fields(line)[fld]
			break
		}
	}
	return result
}
