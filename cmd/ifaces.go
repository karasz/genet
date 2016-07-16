// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of {{ .appName }}, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"net"
	"strings"
)

var ifacesCmd = &cobra.Command{
	Use:   "ifaces",
	Short: "List network interfaces",
	Long:  `List network interfaces and their associated addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		ifaces, err := net.Interfaces()
		if err != nil {
			fmt.Print(fmt.Errorf("Interfaces: %+v\n", err.Error()))
			return
		}
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				fmt.Print(fmt.Errorf("Interfaces: %+v\n", err.Error()))
				continue
			}
			ia := ""
			for l := range addrs {
				ia = ia + "," + addrs[l].String()
			}
			fmt.Printf("%v\t%s \n", i.Name, strings.TrimPrefix(ia, ","))
		}

	},
}

func init() {
	RootCmd.AddCommand(ifacesCmd)
	// TODO: Maybe implement output format like json or so...
	// ifacesCmd.PersistentFlags().String("foo", "", "A help for foo")

}
