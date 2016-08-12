// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package cmd

import (
	"fmt"
	"os"

	"github.com/karasz/genet/genetlib"
	"github.com/spf13/cobra"
)

// iplinkCmd represents the iplink command
var iplinkCmd = &cobra.Command{
	Use:   "iplink",
	Short: "Iplink mimics the ip link group of commands",
	Long:  `Iplink gets informations about link devices using netlink`,
	Run:   doIplink,
}

func init() {
	RootCmd.AddCommand(iplinkCmd)

}

func doIplink(cmd *cobra.Command, args []string) {
	verb := args[0]

	switch verb {
	case "show":
		_, _ = genetlib.LinkShow(len(args) > 1)
	default:
		fmt.Errorf("Unknown or unimplemented verb %s\n", verb)
		os.Exit(1)
	}
}
