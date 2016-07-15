// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of {{ .appName }}, is free and unencumbered
// software released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var GVersion = "0.0.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of genet",
	Long:  `All software has versions. This is genet's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Genet's version is %s\n", GVersion)
	},
}

func init() {
	versionCmd.SetUsageTemplate("Usage: \n\tgenet version\n\n")
	RootCmd.AddCommand(versionCmd)

}
