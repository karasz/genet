package genetlib

import (
	//	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

type TracerouteHop struct {
	Address net.IP
	Host    []string
	RTT     time.Duration
}

type TracerouteResult struct {
	Address net.IP
	Hops    []TracerouteHop
}

func Traceroute(host string) (TracerouteResult, error) {
	var result TracerouteResult

	var dst net.IPAddr
	outIP := GetOutboundIP()
	dst.IP = net.ParseIP(host)
	if dst.IP == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			log.Fatal(err)
		}
		for _, ip := range ips {
			if ip.To4() != nil {
				dst.IP = ip
				result.Address = dst.IP
				break
			}
		}
	} else {
		result.Address = dst.IP
	}

	c, err := net.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	p := ipv4.NewPacketConn(c)

	if err := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		log.Fatal(err)
	}
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}

	rb := make([]byte, 1500)
	for i := 1; i <= 64; i++ {
		var hop TracerouteHop
		wm.Body.(*icmp.Echo).Seq = i
		wb, err := wm.Marshal(nil)

		if err != nil {
			log.Fatal(err)
		}

		if err := p.SetTTL(i); err != nil {
			log.Fatal(err)
		}

		begin := time.Now()

		if _, err := p.WriteTo(wb, nil, &dst); err != nil {
			log.Fatal(err)
		}

		if err := p.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
			log.Fatal(err)
		}

		n, cm, peer, err := p.ReadFrom(rb)

		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			}
			log.Fatal(err)
		}

		rm, err := icmp.ParseMessage(1, rb[:n])

		if err != nil {
			log.Fatal(err)
		}

		if rm.Type == ipv4.ICMPTypeEchoReply {
			if rm.Body.(*icmp.Echo).Seq == i {
				rtt := time.Since(begin)

				names, _ := net.LookupAddr(peer.String())
				hop.Address = net.ParseIP(peer.String())
				hop.Host = names
				hop.RTT = rtt
				result.Hops = append(result.Hops, hop)
				break
			}
		} else {
			if rm.Type == ipv4.ICMPTypeTimeExceeded {
				if cm.Dst.String() == outIP {
					rtt := time.Since(begin)

					names, _ := net.LookupAddr(peer.String())
					hop.Address = net.ParseIP(peer.String())
					hop.Host = names
					hop.RTT = rtt
					result.Hops = append(result.Hops, hop)

				}
			}
		}
	}

	return result, nil
}
