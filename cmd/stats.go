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
	"os"
)

// statsCmd represents the statse command
var statsCmd = &cobra.Command{
	Use:          "stats",
	Short:        "Provide various statistics",
	SilenceUsage: true,
	Long:         `Stats provides various statisctics of the host`,
	Run:          doStats,
}

func init() {
	versionCmd.SetUsageTemplate("Usage: \n\tgenet stats CATEGORY\n\n")
	RootCmd.AddCommand(statsCmd)

}

func doStats(cmd *cobra.Command, args []string) {
	cat := args[0]

	switch cat {
	case "net":
		res, err := genetlib.GetStatsNet()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	case "sys":
		res, err := genetlib.GetStatsSys()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	case "disk":
		res, err := genetlib.GetStatsDisk()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	case "all":
		res, err := genetlib.GetStatsAll()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	default:
		fmt.Errorf("Unknown or unimplemented category %s\n", cat)
		os.Exit(1)

	}

}
