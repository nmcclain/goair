package commands

import (
	"fmt"
	"log"

	"github.com/emccode/govcloudair/types/vcav1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
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
	computeCmd.Flags().StringVar(&region, "region", "", "VCLOUDAIR_REGION")
	computeuseCmd.Flags().StringVar(&planID, "planid", "", "VCLOUDAIR_PLANID")
	computegetCmd.Flags().StringVar(&planID, "planid", "", "VCLOUDAIR_PLANID")
	computeuseCmd.Flags().StringVar(&region, "region", "", "VCLOUDAIR_REGION")
	computegetCmd.Flags().StringVar(&region, "region", "", "VCLOUDAIR_REGION")

	computeCmdV = computeCmd

	computeCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "compute")
		cmd.Usage()
	}
}

var computeuseCmd = &cobra.Command{
	Use:   "use",
	Short: "use",
	Long:  `use compute services`,
	Run:   cmdUseCompute,
}

var computegetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a compute resource from the specified --planid or --region",
	Long:  `Get a compute resource from the specified --planid or --region`,
	Run:   cmdGetCompute,
}

func addCommandsCompute() {
	computeCmd.AddCommand(computeuseCmd)
	computeCmd.AddCommand(computegetCmd)
}

func cmdUseCompute(cmd *cobra.Command, args []string) {
	initConfig(cmd, map[string]FlagValue{
		"planid": {planID, true, false, ""},
		"region": {region, true, false, "planid"},
	})

	err := encodeGobFile("compute", UseValue{
		VarMap: map[string]string{
			"planID": planID,
			"region": region,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if planID != "" {
		fmt.Println(fmt.Sprintf("Set to use PlanID: %v", planID))
	}

	if region != "" {
		fmt.Println(fmt.Sprintf("Set to use region: %v", region))
	}

}

func cmdGetCompute(cmd *cobra.Command, args []string) {
	initConfig(cmd, map[string]FlagValue{})
	client, err := authenticate(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	if err := setGobValues(cmd, "compute"); err != nil {
		log.Fatal(err)
	}

	initConfig(cmd, map[string]FlagValue{
		"planid": {planID, true, false, ""},
	})

	planList, err := client.GetPlans()
	if err != nil {
		log.Fatalf("error Getting plans: %s", err)
	}

	plan := vcatypes.Plan{}
	for _, arg := range planList.Plans {
		if (region != "" && arg.Region == region) || (planID != "" && arg.ID == planID) {
			plan = arg
			break
		}
	}

	yamlOutput, err := yaml.Marshal(&plan)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
