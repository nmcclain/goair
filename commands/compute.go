package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "compute",
	Long:  `compute services`,
}

var computeCmdV *cobra.Command

func init() {
	addCommandsCompute()
	computeCmd.PersistentFlags().StringVar(&username, "username", "", "VCLOUDAIR_USERNAME")
	computeCmd.PersistentFlags().StringVar(&password, "password", "", "VCLOUDAIR_PASSWORD")
	computeCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "VCLOUDAIR_ENDPOINT")
	viper.SetDefault("endpoint", "https://us-california-1-3.vchs.vmware.com/")

	computeCmd.Flags().StringVar(&planID, "planid", "", "VCLOUDAIR_PLANID")
	computeuseCmd.Flags().StringVar(&planID, "planid", "", "VCLOUDAIR_PLANID")

	computeCmdV = computeCmd

	computeCmd.Run = func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		getGobValues(cmd)
	}
}

var computeuseCmd = &cobra.Command{
	Use:   "use",
	Short: "use",
	Long:  `use compute services`,
	Run:   cmdUseCompute,
}

func getGobValues(cmd *cobra.Command) {
	getValue := GetValue{}
	decodeGobFile("compute", &getValue)

	if cmd.Flags().Lookup("planid").Changed == false {
		planID = *getValue.VarMap["planID"]
	}
}

func addCommandsCompute() {
	computeCmd.AddCommand(computeuseCmd)
}

func cmdUseCompute(cmd *cobra.Command, args []string) {
	initConfig(cmd, map[string]FlagValue{
		"planid": {planID, true, false},
	})

	encodeGobFile("compute", UseValue{
		VarMap: map[string]string{
			"planID": planID,
		},
	})
}
