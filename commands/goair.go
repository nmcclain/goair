package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	CfgFile string
)

type FlagValue struct {
	value      string
	mandatory  bool
	persistent bool
}

var GoairCmd = &cobra.Command{
	Use: "goair",
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig()
		cmd.Usage()
	},
}

func Exec() {
	AddCommands()
	GoairCmd.Execute()
}

func AddCommands() {
	GoairCmd.AddCommand(ondemandCmd)
}

var goairCmdV *cobra.Command

func init() {
	GoairCmd.PersistentFlags().StringVar(&CfgFile, "Config", "", "config file (default is $HOME/goair/config.yaml)")
	goairCmdV = GoairCmd
}

func InitConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/goair")
	viper.AddConfigPath("$HOME/.goair/")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VCLOUDAIR")
}
