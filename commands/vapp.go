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

var vappCmdV *cobra.Command

func init() {
	addCommandsVApp()
	vappCmd.Flags().StringVar(&vdchref, "vappname", "", "VCLOUDAIR_VAPPNAME")
	vappCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")
	vappgetCmd.Flags().StringVar(&vdchref, "vappname", "", "VCLOUDAIR_VAPPNAME")
	vappgetCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")

	vappCmdV = vappCmd

	vappCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goair_compute", "")
		cmd.Usage()
	}
}

func addCommandsVApp() {
	vappCmd.AddCommand(vappgetCmd)
}

var vappCmd = &cobra.Command{
	Use:   "vapp",
	Short: "vapp",
	Long:  `vapp`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var vappgetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a vapp",
	Long:  `Get a vapp`,
	Run:   cmdGetVApp,
}

func cmdGetVApp(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
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

	// vdc, err := client.RetrieveVDC()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	vdc := govcloudair.NewVdc(client)
	vdc.Vdc = &types.Vdc{HREF: client.VCDVDCHREF.String()}

	vapps, err := vdc.GetVApp()
	if err != nil {
		log.Fatal(err)
	}

	yamlOutput, err := yaml.Marshal(&vapps)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}
