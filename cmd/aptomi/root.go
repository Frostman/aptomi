package main

import (
	"github.com/Aptomi/aptomi/cmd"
	"github.com/Aptomi/aptomi/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix = "APTOMI"
)

var (
	cfg       = &config.Server{}
	aptomiCmd = &cobra.Command{
		Use:   "aptomi",
		Short: "Aptomi server",
		Long:  "Aptomi server",

		PersistentPreRun: preRun,

		Run: func(cmd *cobra.Command, args []string) {
			// fall back on default help if no args/flags are passed
			cmd.HelpFunc()(cmd, args)
		},
	}
)

func init() {
	viper.SetEnvPrefix(envPrefix)

	cmd.AddDefaultFlags(aptomiCmd, envPrefix)

	cmd.AddStringFlag(aptomiCmd, "db.connection", "db", "", "/etc/aptomi/db.bolt", envPrefix+"_DB_CONN", "DB connection string")

	aptomiCmd.AddCommand(cmd.Version)
}

func preRun(command *cobra.Command, args []string) {
	cmd.ReadConfig(viper.GetViper(), cfg, "/etc/aptomi")
}
