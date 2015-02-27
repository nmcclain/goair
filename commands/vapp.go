package commands

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/emccode/govcloudair"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	types "github.com/vmware/govcloudair/types/v56"
	"gopkg.in/yaml.v1"
)

var vappCmdV *cobra.Command

func init() {
	addCommandsVApp()
	vappCmd.Flags().StringVar(&vappname, "vappname", "", "VCLOUDAIR_VAPPNAME")
	vappCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")
	vappCmd.Flags().StringVar(&vappid, "vappid", "", "VCLOUDAIR_VAPPID")
	vappgetCmd.Flags().StringVar(&vappname, "vappname", "", "VCLOUDAIR_VAPPNAME")
	vappgetCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")
	vappgetCmd.Flags().StringVar(&vappid, "vappid", "", "VCLOUDAIR_VAPPID")
	vappgetstatusCmd.Flags().StringVar(&vappname, "vappname", "", "VCLOUDAIR_VAPPNAME")
	vappgetstatusCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")
	vappgetstatusCmd.Flags().StringVar(&vappid, "vappid", "", "VCLOUDAIR_VAPPID")
	vappactionCmd.Flags().StringVar(&vappname, "vappname", "", "VCLOUDAIR_VAPPNAME")
	vappactionCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")
	vappactionCmd.Flags().StringVar(&vappid, "vappid", "", "VCLOUDAIR_VAPPID")
	vappupdateCmd.Flags().StringVar(&vappname, "vappname", "", "VCLOUDAIR_VAPPNAME")
	vappupdateCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")
	vappupdateCmd.Flags().StringVar(&vappid, "vappid", "", "VCLOUDAIR_VAPPID")
	vappupdateCmd.Flags().StringVar(&cpucount, "cpucount", "", "VCLOUDAIR_CPUCOUNT")
	vappupdateCmd.Flags().StringVar(&memorysizemb, "memorysizemb", "", "VCLOUDAIR_MEMORYSIZEMB")

	vappCmdV = vappCmd

	vappCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goair_compute", "")
		cmd.Usage()
	}
}

func addCommandsVApp() {
	vappCmd.AddCommand(vappgetCmd)
	vappCmd.AddCommand(vappactionCmd)
	vappCmd.AddCommand(vappgetstatusCmd)
	vappCmd.AddCommand(vappupdateCmd)
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

var vappactionCmd = &cobra.Command{
	Use:   "action",
	Short: "action on a vapp",
	Long:  `action on a vapp poweron|poweroff|reboot|reset|suspend|shutdown|undeploy|deploy|delete`,
	Run:   cmdActionVApp,
}

var vappgetstatusCmd = &cobra.Command{
	Use:   "get-status",
	Short: "Get vapp status",
	Long:  `Get vapp status`,
	Run:   cmdGetStatusVApp,
}

var vappupdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update vapp",
	Long:  `Update vapp`,
	Run:   cmdUpdateVApp,
}

func cmdGetVApp(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
		"vappid":             {vappid, false, false, ""},
		"vappname":           {vappname, false, false, "vappid"},
		"instanceAttributes": {instanceAttributes, true, false, ""},
	})

	if len(args) > 1 {
		log.Fatalf("Too many action statements")
	}

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

	vapp := *govcloudair.NewVApp(client)

	if viper.GetString("vappname") != "" {
		vapp, err = vdc.FindVAppByName(viper.GetString("vappname"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if viper.GetString("vappid") != "" {
		vapp, err = vdc.FindVAppByID(viper.GetString("vappid"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(args) == 0 && vapp.VApp != nil && viper.GetString("vappname") == "" && viper.GetString("vappid") == "" {
		yamlOutput, err := yaml.Marshal(&vapp)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	if len(args) == 1 {
		if vapp.VApp == nil {
			log.Fatalf("err vapp not found")
		}
		extra := args[0]

		switch extra {
		case "vm":
			vm := vapp.VApp.Children.VM
			yamlOutput, err := yaml.Marshal(vm)
			if err != nil {
				log.Fatalf("error marshaling: %s", err)
			}
			fmt.Println(string(yamlOutput))
			return
		case "guestcustomization":
			guestCustomizationSection, err := vapp.GetGuestCustomization()
			if err != nil {
				log.Fatalf("error getting guest customization: %s", err)
			}

			yamlOutput, err := yaml.Marshal(&guestCustomizationSection)
			fmt.Println(string(yamlOutput))
			return
		default:
			log.Fatalf("Did not specify proper extra get command of vm|guestcustomization")
		}
	}

	vapps, err := vdc.GetVApp()
	if err != nil {
		log.Fatal(err)
	}
	yamlOutput, err := yaml.Marshal(&vapps)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))
	return

}

func cmdActionVApp(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
		"vappid":             {vappid, true, false, ""},
		"vappname":           {vappname, true, false, "vappid"},
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

	vapp := *govcloudair.NewVApp(client)
	if viper.GetString("vappname") != "" {
		vapp, err = vdc.FindVAppByName(viper.GetString("vappname"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if viper.GetString("vappid") != "" {
		vapp, err = vdc.FindVAppByID(viper.GetString("vappid"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if vapp.VApp == nil {
		log.Fatalf("Couldn't find VApp")
	}

	if len(args) > 1 {
		log.Fatalf("Too many action statements")
	}

	action := args[0]

	switch action {
	case "poweron":
		vapp.PowerOn()
	case "poweroff":
		vapp.PowerOff()
	case "reboot":
		vapp.Reboot()
	case "reset":
		vapp.Reset()
	case "suspend":
		vapp.Suspend()
	case "shutdown":
		vapp.Shutdown()
	case "undeploy":
		vapp.Undeploy()
	case "deploy":
		vapp.Deploy()
	case "delete":
		vapp.Delete()
	}

	return

}

func cmdGetStatusVApp(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
		"vappid":             {vappid, true, false, ""},
		"vappname":           {vappname, true, false, "vappid"},
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

	vapp := *govcloudair.NewVApp(client)
	if viper.GetString("vappname") != "" {
		vapp, err = vdc.FindVAppByName(viper.GetString("vappname"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if viper.GetString("vappid") != "" {
		vapp, err = vdc.FindVAppByID(viper.GetString("vappid"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if vapp.VApp == nil {
		log.Fatalf("Couldn't find VApp")
	}

	vappStatus, err := vapp.GetStatus()
	if err != nil {
		log.Fatalf("error getting VApp status: %v", err)
	}

	fmt.Println(vappStatus)
}

func cmdUpdateVApp(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
		"vappid":             {vappid, true, false, ""},
		"vappname":           {vappname, true, false, "vappid"},
		"memorysizemb":       {memorysizemb, false, false, "memorysizemb"},
		"cpucount":           {cpucount, false, false, "cpucount"},
		"runasync":           {cpucount, false, false, "runasync"},
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

	vapp := *govcloudair.NewVApp(client)
	if viper.GetString("vappname") != "" {
		vapp, err = vdc.FindVAppByName(viper.GetString("vappname"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if viper.GetString("vappid") != "" {
		vapp, err = vdc.FindVAppByID(viper.GetString("vappid"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if vapp.VApp == nil {
		log.Fatalf("Couldn't find VApp")
	}

	if memorysizemb != "" {
		ms, _ := strconv.Atoi(memorysizemb)
		task, err := vapp.ChangeMemorySize(ms)
		if err != nil {
			log.Fatalf("err problem changing memory size: %v", err)
		}

		err = task.WaitTaskCompletion()
		if err != nil {
			log.Fatalf("error waiting for task to complete: %v", err)
		}
	}

	if cpucount != "" {
		cc, _ := strconv.Atoi(cpucount)
		task, err := vapp.ChangeCPUcount(cc)
		if err != nil {
			log.Fatalf("err problem updating cpu count: %v", err)
		}

		err = task.WaitTaskCompletion()
		if err != nil {
			log.Fatalf("error waiting for task to complete: %v", err)
		}
	}

}
