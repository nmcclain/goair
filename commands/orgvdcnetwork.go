package commands

import (
	"fmt"
	"log"
	"net/url"

	"github.com/emccode/govcloudair"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	types "github.com/vmware/govcloudair/types/v56"
	"gopkg.in/yaml.v1"
)

var orgvdcnetworkCmdV *cobra.Command

func init() {
	addCommandsOrgVdcNetwork()
	orgvdcnetworkCmd.Flags().StringVar(&vdcnetworkname, "vdcnetworkname", "", "VCLOUDAIR_VDCNETWORKNAME")
	orgvdcnetworkgetCmd.Flags().StringVar(&vdcnetworkname, "vdcnetworkname", "", "VCLOUDAIR_VDCNETWORKNAME")

	orgvdcnetworkCmdV = orgvdcnetworkCmd

	orgvdcnetworkCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goair_compute", "")
		cmd.Usage()
	}
}

func addCommandsOrgVdcNetwork() {
	orgvdcnetworkCmd.AddCommand(orgvdcnetworkgetCmd)
}

var orgvdcnetworkCmd = &cobra.Command{
	Use:   "orgvdcnetwork",
	Short: "orgvdcnetwork",
	Long:  `orgvdcnetwork`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var orgvdcnetworkgetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get orgvdcnetworks or specific network with --vdcnetworkname",
	Long:  `Get orgvdcnetworks or specific network with --vdcnetworkname`,
	Run:   cmdGetOrgVdcNetwork,
}

func cmdGetOrgVdcNetwork(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
		"vdcnetworkname":     {vdcnetworkname, false, false, ""},
		"instanceAttributes": {instanceAttributes, true, false, ""},
	})

	client, err := authenticate(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	err = authenticatecompute(client, false, "")
	if err != nil {
		log.Fatalf("Error authenticating compute: %s", err)
	}

	vdcuri, err := url.Parse(viper.GetString("vdchref"))
	if err != nil {
		log.Fatal(err)
	}

	client.VCDVDCHREF = *vdcuri

	vdc := govcloudair.NewVdc(client)
	vdc.Vdc = &types.Vdc{HREF: client.VCDVDCHREF.String()}
	err = vdc.Refresh()
	if err != nil {
		log.Fatalf("err refreshing vdc: %v", err)
	}

	if viper.GetString("vdcnetworkname") != "" {
		VdcNetwork, err := vdc.FindVDCNetwork(viper.GetString("vdcnetworkname"))
		if err != nil {
			log.Fatalf("error finding Vdc network: %v", err)
		}

		yamlOutput, err := yaml.Marshal(&VdcNetwork)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	VdcNetworks, err := vdc.GetVDCNetwork()
	if err != nil {
		log.Fatalf("error finding Vdc network: %v", err)
	}

	yamlOutput, err := yaml.Marshal(&VdcNetworks)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))
	return

}
