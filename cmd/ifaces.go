// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package cmd

import (
	"fmt"
	"github.com/karasz/genet/genetlib"
	"github.com/spf13/cobra"
)

var ifacesCmd = &cobra.Command{
	Use:   "ifaces",
	Short: "List network interfaces",
	Long:  `List network interfaces and their associated addresses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ifaces, err := genetlib.GetIfaces()
		if err != nil {
			fmt.Printf("The following error occured: %s\n", err)
			return err
		}
		fmt.Printf("%+v\n", ifaces)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(ifacesCmd)
	// TODO: Maybe implement output format like json or so...
	// ifacesCmd.PersistentFlags().String("foo", "", "A help for foo")

}
