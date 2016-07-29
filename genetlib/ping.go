// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/karasz/go-fastping"
)

type pingResp struct {
	addr *net.IPAddr
	rtt  time.Duration
}

type pingStats struct {
	Responses []pingResp
	Sent      int
	Received  int
	Lost      int
	Totaltime time.Duration
	Min       time.Duration
	Max       time.Duration
	Avg       time.Duration
}

func (p pingStats) String() string {
	return fmt.Sprintf("Sent: %d, Received: %d, Lost: %d, Total time: %v, Min: %v, Max: %v, Avg: %v", p.Sent, p.Received, p.Lost, p.Totaltime, p.Min, p.Max, p.Avg)
}

func Ping(addr string, prot string, cnt int, iface string, statsonly bool) (pingStats, error) {

	if strings.ToLower(prot) == "tcp" {
		return pingTCP(addr, cnt, iface, statsonly)
	}

	return pingICMP(addr, prot, cnt, iface, statsonly)
}

func pingICMP(address string, prot string, cnt int, iface string, statsonly bool) (pingStats, error) {
	var stats pingStats
	var err error

	min := time.Hour
	max := time.Nanosecond
	alltime := time.Nanosecond * 0

	p := fastping.NewPinger()

	if strings.ToLower(prot) == "udp" {
		p.Network("udp")
	}
	netProto := "ip4:icmp"

	if strings.Index(address, ":") != -1 {
		netProto = "ip6:ipv6-icmp"
	}

	ra, err := net.ResolveIPAddr(netProto, address)

	if err != nil {
		fmt.Errorf("The following error occured %s", err)
		os.Exit(1)
	}

	if iface != "" {
		p.Source(iface)
	}

	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		resp := pingResp{addr, t}
		alltime = alltime + t
		if t < min {
			min, stats.Min = t, t
		}

		if t > max {
			max, stats.Max = t, t
		}
		stats.Responses = append(stats.Responses, resp)
		stats.Received = stats.Received + 1
	}

	p.MaxRTT = time.Second
	start := time.Now()

	for i := 1; i <= cnt; i++ {
		if !statsonly {
			fmt.Print(". ")
		}
		p.Run()
		stats.Sent = stats.Sent + 1
	}
	if !statsonly {
		fmt.Println("")
	}
	stats.Avg = alltime / time.Duration(cnt)
	stats.Totaltime = time.Since(start)
	stats.Lost = stats.Sent - stats.Received

	return stats, err
}

func pingTCP(addr string, cnt int, iface string, q bool) (pingStats, error) {
	var stats pingStats
	err := errors.New("TCP ping not implemented")
	return stats, err
}
