package cmd

import (
	"fmt"
	"github.com/Frostman/aptomi/pkg/slinga"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show an object",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var showCmdConfig = &cobra.Command{
	Use:   "config",
	Short: "Show aptomi configuration variables",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		vars := []string{"APTOMI_POLICY", "APTOMI_DB"}

		for _, key := range vars {
			value, _ := os.LookupEnv(key)
			fmt.Println(key + " = " + value)
		}
	},
}

var showCmdPolicy = &cobra.Command{
	Use:   "policy",
	Short: "Show aptomi policy",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var showCmdGraph = &cobra.Command{
	Use:   "graph",
	Short: "Show the current allocation graph (what has been allocated and who is using what)",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		usage := slinga.LoadServiceUsageState()
		usage.DrawVisualAndStore()

		command := exec.Command("open", []string{usage.GetVisualFileNamePNG()}...)
		if err := command.Run(); err != nil {
			fmt.Print("Allocations (PNG): " + usage.GetVisualFileNamePNG())
		}
	},
}

func init() {
	showCmd.AddCommand(showCmdConfig)
	showCmd.AddCommand(showCmdPolicy)
	showCmd.AddCommand(showCmdGraph)

	RootCmd.AddCommand(showCmd)
}
