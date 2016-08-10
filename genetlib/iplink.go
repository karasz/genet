// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"fmt"
	"github.com/vishvananda/netlink"
	//	"net"
	"strings"
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
	result.state = "<NOT IMPLEMENTED>"
	result.mode = "<NOT IMPLEMENTED>"
	result.group = "<NOT IMPLEMENTED>"
	result.qlen = 0

	return result
}

func LinkShow() ([]netlink.Link, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return links, err
	}
	for _, link := range links {
		z := fillAttrs(link)
		fmt.Printf("%d: %s: <%s> mtu %d qdisc %s state %s mode %s group %s qlen %d\n\tlink/%s: %s brd: %s\n",
			z.index, z.name, z.flags, z.mtu, z.qdisc, z.state, z.mode, z.group, z.qlen, z.linktype, z.hwd, z.brd)
	}
	return links, nil
}
