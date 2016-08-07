// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of {{ .appName }}, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package cmd

import (
	"fmt"
	"github.com/karasz/genet/genetlib"
	"github.com/spf13/cobra"
	"os"
)

// tracerouteCmd represents the traceroute command
var tracerouteCmd = &cobra.Command{
	Use:          "traceroute",
	Short:        "Trace route to host",
	SilenceUsage: true,
	Long:         `Traceroute emulates the traceroute command from linux`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must provide a host to ping as the first argument")
		}
		return nil
	},
	Run: doTraceroute,
}

func init() {
	versionCmd.SetUsageTemplate("Usage: \n\tgenet traceroute HOST\n\n")
	RootCmd.AddCommand(tracerouteCmd)

}

func doTraceroute(cmd *cobra.Command, args []string) {
	addr := args[0]
	if addr == "" {
		fmt.Errorf("An empty address was provided. You must provide an address as the first argument.\n")
		os.Exit(1)
	}

	res, err := genetlib.Traceroute(addr)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
