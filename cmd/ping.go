// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package cmd

import (
	"errors"
	"fmt"
	"github.com/karasz/genet/genetlib"
	"github.com/spf13/cobra"
	"os"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:          "ping HOST",
	Short:        "Send packets to network hosts",
	Long:         `Ping emulates a ping command sending packets to network hosts`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must provide a host to ping as the first argument")
		}
		return nil
	},
	Run: doPing,
}

func init() {
	RootCmd.AddCommand(pingCmd)

	pingCmd.Flags().StringP("protocol", "p", "ICMP", "Protocol to use for ping")
	pingCmd.Flags().StringP("interface", "I", "", "Interface to use for ping")
	pingCmd.Flags().IntP("count", "c", 10, "Number of packets to send")

}

func doPing(cmd *cobra.Command, args []string) {
	addr := args[0]
	if addr == "" {
		fmt.Fprintf(os.Stderr, "An empty address was provided. You must provide an address as the first argument.\n")
		os.Exit(1)
	}
	prot, _ := cmd.Flags().GetString("protocol")
	iface, _ := cmd.Flags().GetString("interface")
	cnt, _ := cmd.Flags().GetInt("count")

	res, err := genetlib.Ping(addr, prot, cnt, iface)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
