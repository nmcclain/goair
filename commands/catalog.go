package commands

import (
	"fmt"
	"log"
	"net/url"

	"github.com/emccode/govcloudair"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

var catalogCmdV *cobra.Command

func init() {
	addCommandsCatalog()
	catalogCmd.Flags().StringVar(&catalogname, "catalogname", "", "VCLOUDAIR_CATALOGNAME")
	cataloggetCmd.Flags().StringVar(&catalogname, "catalogname", "", "VCLOUDAIR_CATALOGNAME")

	catalogCmdV = catalogCmd

	catalogCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goair_compute", "")
		cmd.Usage()
	}
}

func addCommandsCatalog() {
	catalogCmd.AddCommand(cataloggetCmd)
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

func cmdGetCatalog(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":      {planID, true, false, ""},
		"region":      {region, true, false, "planid"},
		"orghref":     {orghref, true, false, ""},
		"catalogname": {catalogname, true, false, ""},
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

	catalogItems, err := catalog.GetCatalogItem()

	yamlOutput, err := yaml.Marshal(&catalogItems)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))
	return

}
