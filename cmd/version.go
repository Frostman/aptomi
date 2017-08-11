package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var Version string
var BuildTime string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version string of Aptomi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Aptomi commit: %s\n       built: %s\n", Version, BuildTime)
	},
}