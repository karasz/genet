// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/vishvananda/netlink"
)

type genetlink struct {
	index    int
	name     string
	linktype string
	flags    string
	mtu      int
	qdisc    string
	state    string
	mode     string
	group    string
	qlen     int
	hwd      string
	brd      string
}

type linkstats struct {
	rxPackets         uint32
	txPackets         uint32
	rxBytes           uint32
	txBytes           uint32
	rxErrors          uint32
	txErrors          uint32
	rxDropped         uint32
	txDropped         uint32
	multicast         uint32
	collisions        uint32
	rxLengthErrors    uint32
	rxOverErrors      uint32
	rxCrcErrors       uint32
	rxFrameErrors     uint32
	rxFifoErrors      uint32
	rxMissedErrors    uint32
	txAbortedErrors   uint32
	txCarrierErrors   uint32
	txFifoErrors      uint32
	txHeartbeatErrors uint32
	txWindowErrors    uint32
	rxCompressed      uint32
	txCompressed      uint32
}

func isLoopback(z netlink.Link) bool {
	linka := z.Attrs()
	// 2 is loopback
	if linka.Flags&(1<<uint(2)) != 0 {
		return true
	}
	return false
}

func fillAttrs(z netlink.Link) genetlink {
	var result genetlink
	linka := z.Attrs()
	result.index = linka.Index
	result.name = linka.Name
	result.flags = strings.Replace(strings.ToUpper(linka.Flags.String()), "|", ",", -1)
	result.mtu = linka.MTU
	result.qlen = linka.TxQLen

	qds, _ := netlink.QdiscList(z)
	for _, qd := range qds {
		qda := qd.Attrs()
		if qda.Parent == netlink.HANDLE_ROOT {
			result.qdisc = qd.Type()
			break
		}
	}

	if linka.HardwareAddr != nil {
		result.hwd = linka.HardwareAddr.String()
		result.brd = "ff:ff:ff:ff:ff:ff"
	} else {
		result.hwd = "00:00:00:00:00:00"
		result.brd = "00:00:00:00:00:00"
	}
	// just for similarity with the iproute2 output
	if isLoopback(z) {
		result.linktype = "loopback"
	} else {
		if z.Type() != "device" {
			result.linktype = z.Type()
		} else {
			result.linktype = "ether"
		}
	}

	result.state = GetValFromFile("/sys/class/net/"+result.name+"/operstate", 0)
	if GetValFromFile("/sys/class/net/"+result.name+"/link_mode", 0) == "0" {
		result.mode = "DEFAULT"
	} else {
		result.mode = "DORMANT"
	}
	result.group = GetFieldValFromFile("/etc/iproute2/group", 0, GetValFromFile("/sys/class/net/"+result.name+"/netdev_group", 0), 1)

	return result
}

func fillStats(z netlink.Link) linkstats {
	var result linkstats
	mst := z.Attrs().Statistics

	result.rxPackets = mst.RxPackets
	result.txPackets = mst.TxPackets
	result.rxBytes = mst.RxBytes
	result.txBytes = mst.TxBytes
	result.rxErrors = mst.RxErrors
	result.txErrors = mst.TxErrors
	result.rxDropped = mst.RxDropped
	result.txDropped = mst.TxDropped
	result.multicast = mst.Multicast
	result.collisions = mst.Collisions
	result.rxLengthErrors = mst.RxLengthErrors
	result.rxOverErrors = mst.RxOverErrors
	result.rxCrcErrors = mst.RxCrcErrors
	result.rxFrameErrors = mst.RxFrameErrors
	result.rxFifoErrors = mst.RxFifoErrors
	result.rxMissedErrors = mst.RxMissedErrors
	result.txAbortedErrors = mst.TxAbortedErrors
	result.txCarrierErrors = mst.TxCarrierErrors
	result.txFifoErrors = mst.TxFifoErrors
	result.txHeartbeatErrors = mst.TxHeartbeatErrors
	result.txWindowErrors = mst.TxWindowErrors
	result.rxCompressed = mst.RxCompressed
	result.txCompressed = mst.TxCompressed

	return result
}

func LinkShow(detailed bool) ([]netlink.Link, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return links, err
	}
	for _, link := range links {
		z := fillAttrs(link)
		fmt.Printf("%d: %s: <%s> mtu %d qdisc %s state %s mode %s group %s qlen %d\n\tlink/%s: %s brd: %s\n",
			z.index, z.name, z.flags, z.mtu, z.qdisc, z.state, z.mode, z.group, z.qlen, z.linktype, z.hwd, z.brd)
		if detailed {
			d := fillStats(link)
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, " \tRX: bytes\tpackets\terrors\tdropped\toverrun\tmcast")
			fmt.Fprintf(w, " \t%d\t%d\t%d\t%d\t%d\t%d\n", d.rxBytes, d.rxPackets, d.rxErrors, d.rxDropped, d.rxOverErrors, d.multicast)
			fmt.Fprintln(w, " \tRX errors:\tlength\tcrc\tframe\tfifo\tmissed")
			fmt.Fprintf(w, " \t%d\t%d\t%d\t%d\t%d\n", d.rxLengthErrors, d.rxCrcErrors, d.rxFrameErrors, d.rxFifoErrors, d.rxMissedErrors)
			fmt.Fprintln(w, " \tTX: bytes\tpackets\terrors\tdropped\tcarrier\tcollsns")
			fmt.Fprintf(w, " \t%d\t%d\t%d\t%d\t%d\t%d\n", d.txBytes, d.txPackets, d.txErrors, d.txDropped, d.txCarrierErrors, d.collisions)
			fmt.Fprintln(w, " \tTX errors:\taborted\tfifo\twindow\theartbeat")
			fmt.Fprintf(w, " \t%d\t%d\t%d\t%d\n", d.txAbortedErrors, d.txFifoErrors, d.txWindowErrors, d.txHeartbeatErrors)
			w.Flush()
		}
	}
	return links, nil
}
