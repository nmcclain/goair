package commands

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/emccode/clue"
	"github.com/emccode/govcloudair"
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
		setGobValues(cmd, "goair_compute", "")
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

func authenticatecompute(client *govcloudair.ODClient, force bool, ia string) (err error) {
	getValue := clue.GetValue{}
	if err := clue.DecodeGobFile("goair_compute", &getValue); err != nil {
		return fmt.Errorf("Problem with client DecodeGobFile", err)
	}

	if ia != "" {
		if getValue.VarMap["instanceAttributes"] != nil {
			ia = *getValue.VarMap["instanceAttributes"]
		}
	}

	if force || *getValue.VarMap["VCDToken"] == "" {

		instanceAttributes := vcatypes.InstanceAttributes{}
		json.Unmarshal([]byte(ia), &instanceAttributes)

		err = client.GetBackendAuth(instanceAttributes)
		if err != nil {
			return fmt.Errorf("error Authenticating: %s", err)
		}

		var planid string
		var region string

		if getValue.VarMap["planID"] != nil {
			planid = *getValue.VarMap["planID"]
		}

		if getValue.VarMap["region"] != nil {
			region = *getValue.VarMap["region"]
		}

		err = clue.EncodeGobFile("goair_compute", clue.UseValue{
			VarMap: map[string]string{
				"planID":             planid,
				"region":             region,
				"instanceAttributes": ia,
				"VCDToken":           client.VCDToken,
			},
		})
		if err != nil {
			return fmt.Errorf("Error encoding gob: %s", err)
		}

	} else {
		client.VCDToken = *getValue.VarMap["VCDToken"]
	}

	return nil
}

func cmdUseCompute(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", false, map[string]FlagValue{
		"planid": {planID, true, false, ""},
		"region": {region, true, false, "planid"},
	})

	client, err := authenticate(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	instanceList, err := client.GetInstances()
	if err != nil {
		log.Fatalf("error Getting instances: %s", err)
	}

	instance := govcloudair.Instance{}
	for _, arg := range instanceList.Instances {
		if (viper.GetString("region") != "" && arg.Region == region) || (viper.GetString("planid") != "" && arg.PlanID == planID) {
			instance = govcloudair.Instance(arg)
			break
		}
	}

	err = authenticatecompute(client, true, instance.InstanceAttributes)
	if err != nil {
		log.Fatalf("Error authenticating compute: %s", err)
	}

	if planID != "" {
		fmt.Println(fmt.Sprintf("Set to use PlanID: %v", planID))
	}

	if region != "" {
		fmt.Println(fmt.Sprintf("Set to use region: %v", region))
	}

}

func cmdGetCompute(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid": {planID, true, false, ""},
		"region": {region, true, false, "planid"},
	})

	client, err := authenticate(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	instanceList, err := client.GetInstances()
	if err != nil {
		log.Fatalf("error Getting instances: %s", err)
	}

	instance := govcloudair.Instance{}
	for _, arg := range instanceList.Instances {
		if (viper.GetString("region") != "" && arg.Region == region) || (viper.GetString("planid") != "" && arg.PlanID == planID) {
			instance = govcloudair.Instance(arg)
			break
		}
	}

	err = authenticatecompute(client, true, instance.InstanceAttributes)
	if err != nil {
		fmt.Errorf("Errgeting authenticating compute: %s", err)
	}
	links, err := govcloudair.GetOrgVdc(client, &client.VCDORGHREF)

	yamlOutput, err := yaml.Marshal(&links)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
