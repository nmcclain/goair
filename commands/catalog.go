package commands

import (
	"fmt"
	"log"
	"net/url"

	"github.com/emccode/govcloudair"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmware/govcloudair/types/v56"
	"gopkg.in/yaml.v1"
)

var catalogCmdV *cobra.Command

func init() {
	addCommandsCatalog()
	catalogCmd.Flags().StringVar(&catalogname, "catalogname", "", "VCLOUDAIR_CATALOGNAME")
	catalogCmd.Flags().StringVar(&catalogitemname, "catalogitemname", "", "VCLOUDAIR_CATALOGITEMNAME")
	cataloggetCmd.Flags().StringVar(&catalogname, "catalogname", "", "VCLOUDAIR_CATALOGNAME")
	cataloggetCmd.Flags().StringVar(&catalogitemname, "catalogitemname", "", "VCLOUDAIR_CATALOGITEMNAME")
	catalogdeployCmd.Flags().StringVar(&catalogname, "catalogname", "", "VCLOUDAIR_CATALOGNAME")
	catalogdeployCmd.Flags().StringVar(&catalogitemname, "catalogitemname", "", "VCLOUDAIR_CATALOGITEMNAME")
	catalogdeployCmd.Flags().StringVar(&vmname, "vmname", "", "VCLOUDAIR_VMNAME")
	catalogdeployCmd.Flags().StringVar(&vdcnetworkname, "vdcnetworkname", "", "VCLOUDAIR_VDCNETWORKNAME")
	catalogdeployCmd.Flags().StringVar(&runasync, "runasync", "", "VCLOUDAIR_RUNASYNC")

	catalogCmdV = catalogCmd

	catalogCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goair_compute", "")
		cmd.Usage()
	}
}

func addCommandsCatalog() {
	catalogCmd.AddCommand(cataloggetCmd)
	catalogCmd.AddCommand(catalogdeployCmd)
}

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "catalog",
	Long:  `catalog`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var cataloggetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a catalog",
	Long:  `Get a catalog`,
	Run:   cmdGetCatalog,
}

var catalogdeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy from catalog item",
	Long:  `Deploy from catalog item`,
	Run:   cmdDeployCatalog,
}

func cmdGetCatalog(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":          {planID, true, false, ""},
		"region":          {region, true, false, "planid"},
		"orghref":         {orghref, true, false, ""},
		"catalogname":     {catalogname, false, false, ""},
		"catalogitemname": {catalogitemname, false, false, ""},
	})

	client, err := authenticate(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	err = authenticatecompute(client, false, "")
	if err != nil {
		log.Fatalf("Error authenticating compute: %s", err)
	}

	orgURI, err := url.ParseRequestURI(viper.GetString("orghref"))
	if err != nil {
		log.Fatalf("err parsing org uri: %v", err)
	}

	client.VCDORGHREF = *orgURI

	org, err := govcloudair.GetOrg(client, &client.VCDORGHREF)
	if err != nil {
		log.Fatalf("Error getting org: %v", err)
	}

	if viper.GetString("catalogname") == "" {
		catalogs, err := org.GetCatalog()
		if err != nil {
			log.Fatalf("Error getting catalogs: %v", err)
		}

		yamlOutput, err := yaml.Marshal(&catalogs)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	catalog, err := org.FindCatalog(catalogname)
	if err != nil {
		log.Fatalf("err: problem finding catalog: %v", err)
	}

	if viper.GetString("catalogitemname") == "" {
		catalogItems, err := catalog.GetCatalogItem()

		yamlOutput, err := yaml.Marshal(&catalogItems)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	catalogItem, err := catalog.FindCatalogItem(catalogitemname)
	if err != nil {
		log.Fatalf("err: problem finding catalog: %v", err)
	}

	if args[0] == "" {
		yamlOutput, err := yaml.Marshal(&catalogItem)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	if args[0] == "vapptemplate" {
		vappTemplate, err := catalogItem.GetVAppTemplate()
		if err != nil {
			log.Fatalf("err: problem getting VApp Template: %v", err)
		}

		yamlOutput, err := yaml.Marshal(&vappTemplate)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	log.Fatalf("Problem with get action specified: %v", args[0])

}

func cmdDeployCatalog(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":          {planID, true, false, ""},
		"region":          {region, true, false, "planid"},
		"orghref":         {orghref, true, false, ""},
		"vdchref":         {vdchref, true, false, ""},
		"catalogname":     {catalogname, true, false, ""},
		"catalogitemname": {catalogitemname, true, false, ""},
		"vmname":          {vmname, true, false, ""},
		"vdcnetworkname":  {vdcnetworkname, true, false, ""},
	})

	client, err := authenticate(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	err = authenticatecompute(client, false, "")
	if err != nil {
		log.Fatalf("Error authenticating compute: %s", err)
	}

	orgURI, err := url.ParseRequestURI(viper.GetString("orghref"))
	if err != nil {
		log.Fatalf("err parsing org uri: %v", err)
	}

	client.VCDORGHREF = *orgURI

	org, err := govcloudair.GetOrg(client, &client.VCDORGHREF)
	if err != nil {
		log.Fatalf("Error getting org: %v", err)
	}

	catalog, err := org.FindCatalog(catalogname)
	if err != nil {
		log.Fatalf("err: problem finding catalog: %v", err)
	}

	catalogItem, err := catalog.FindCatalogItem(catalogitemname)
	if err != nil {
		log.Fatalf("err: problem finding catalog: %v", err)
	}

	vappTemplate, err := catalogItem.GetVAppTemplate()
	if err != nil {
		log.Fatalf("err: problem getting VApp Template: %v", err)
	}

	vdcuri, err := url.Parse(viper.GetString("vdchref"))
	if err != nil {
		log.Fatal(err)
	}

	client.VCDVDCHREF = *vdcuri

	vdc := govcloudair.NewVdc(client)
	vdc.Vdc = &types.Vdc{HREF: client.VCDVDCHREF.String()}
	vdc.Refresh()

	if _, err := vdc.FindVAppByName(vappname); err == nil {
		log.Fatalf("VApp %v already exists", viper.GetString("vappname"))
	}

	vdcNetwork, err := vdc.FindVDCNetwork(viper.GetString("vdcnetworkname"))
	if err != nil {
		log.Fatalf("error finding Vdc network: %v", err)
	}

	vapp := govcloudair.NewVApp(client)
	task, err := vapp.ComposeVApp(vdcNetwork, vappTemplate, "", "", viper.GetString("vmname"))
	if err != nil {
		log.Fatalf("error composing vapp: %v", err)
	}

	if viper.GetString("runasync") == "true" {
		yamlOutput, err := yaml.Marshal(&task)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	err = task.WaitTaskCompletion()
	if err != nil {
		log.Fatalf("error waiting for task to complete: %v", err)
	}

}
