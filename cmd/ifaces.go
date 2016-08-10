// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var ifacesCmd = &cobra.Command{
	Use:   "ifaces",
	Short: "List network interfaces",
	Long:  `List network interfaces and their associated addresses.`,
	RunE:  doIfaces,
}

func init() {
	RootCmd.AddCommand(ifacesCmd)
}

func doIfaces(cmd *cobra.Command, args []string) error {
	result := ""
	ifaces, err := ioutil.ReadDir("/sys/class/net/")
	if err != nil {
		fmt.Println(result)
		return err
	}

	for _, iface := range ifaces {
		fmt.Println(iface.Name())
	}

	return nil
}
