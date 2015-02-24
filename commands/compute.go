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
	types "github.com/vmware/govcloudair/types/v56"
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
	computeuseCmd.Flags().StringVar(&vdcname, "vdcname", "", "VCLOUDAIR_VDCNAME")
	computeuseCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")

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

func authenticatecompute(client *govcloudair.Client, force bool, ia string) (err error) {
	getValue := clue.GetValue{}
	if err := clue.DecodeGobFile("goair_compute", &getValue); err != nil {
		return fmt.Errorf("Problem with client DecodeGobFile: %v", err)
	}

	if ia != "" {
		if getValue.VarMap["instanceAttributes"] != nil {
			ia = *getValue.VarMap["instanceAttributes"]
		}
	}

	if force || *getValue.VarMap["VCDToken"] == "" || *getValue.VarMap["VCDAuthHeader"] == "" {

		instanceAttributes := vcatypes.InstanceAttributes{}
		json.Unmarshal([]byte(ia), &instanceAttributes)

		err = client.GetBackendAuthOD(instanceAttributes)
		if err != nil {
			return fmt.Errorf("error Authenticating: %s", err)
		}

		err = clue.EncodeGobFile("goair_compute", clue.UseValue{
			VarMap: map[string]string{
				"planID":             viper.GetString("planid"),
				"region":             viper.GetString("region"),
				"instanceAttributes": ia,
				"VCDToken":           client.VCDToken,
				"VCDAuthHeader":      client.VCDAuthHeader,
			},
		})
		if err != nil {
			return fmt.Errorf("Error encoding gob: %s", err)
		}

	} else {
		client.VCDToken = *getValue.VarMap["VCDToken"]
		client.VCDAuthHeader = *getValue.VarMap["VCDAuthHeader"]
	}

	return nil
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
		if (viper.GetString("region") != "" && arg.Region == viper.GetString("region")) || (viper.GetString("planid") != "" && arg.PlanID == viper.GetString("planid")) {
			instance = govcloudair.Instance(arg)
			break
		}
		log.Fatalf("Couldn't find region or planid")
	}

	err = authenticatecompute(client, true, instance.InstanceAttributes)
	if err != nil {
		log.Fatalf("Err authenticating compute: %s", err)
	}

	links, err := govcloudair.GetOrgVdc(client, &client.VCDORGHREF)
	if err != nil {
		log.Fatalf("Err geting orgvdc: %s", err)
	}

	yamlOutput, err := yaml.Marshal(&links)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}

func cmdUseCompute(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":  {planID, true, false, ""},
		"region":  {region, true, false, "planid"},
		"vdchref": {vdchref, true, false, ""},
		"vdcname": {vdcname, true, false, "vdchref"},
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

	links, err := govcloudair.GetOrgVdc(client, &client.VCDORGHREF)

	useLink := types.Link{}
	for _, link := range links {
		//fmt.Printf("vdchref: %v\nvdcname: %v\n", viper.GetString("vdchref"), viper.GetString("vdcname"))
		if (viper.GetString("vdchref") != "" && link.HREF == viper.GetString("vdchref")) ||
			(viper.GetString("vdcname") != "" && link.Name == viper.GetString("vdcname")) {
			useLink = link
			break
		}
	}

	if useLink.Name != "" {
		err = clue.EncodeGobFile("goair_compute", clue.UseValue{
			VarMap: map[string]string{
				"planID":             viper.GetString("planid"),
				"region":             viper.GetString("region"),
				"vdchref":            useLink.HREF,
				"instanceAttributes": instance.InstanceAttributes,
				"VCDToken":           client.VCDToken,
				"VCDAuthHeader":      client.VCDAuthHeader,
			},
		})

		yamlOutput, err := yaml.Marshal(&useLink)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))

	} else {
		log.Fatalf("Failed to find VDC %v%v", vdcname, vdchref)
	}
}
