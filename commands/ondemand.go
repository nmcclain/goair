package commands

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/emccode/goair/table"
	"github.com/emccode/govcloudair"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	username string
	password string
	endpoint string
)

var ondemandCmd = &cobra.Command{
	Use:   "ondemand",
	Short: "ondemand",
	Long:  `Get plans`,
}

var ondemandCmdV *cobra.Command

func init() {
	addCommands()
	ondemandCmd.PersistentFlags().StringVar(&username, "username", "", "VCLOUDAIR_USERNAME")
	ondemandCmd.PersistentFlags().StringVar(&password, "password", "", "VCLOUDAIR_PASSWORD")
	ondemandCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "VCLOUDAIR_ENDPOINT")
	viper.SetDefault("endpoint", "https://us-california-1-3.vchs.vmware.com/")

	ondemandCmdV = ondemandCmd

	ondemandCmd.Run = func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	}

}

func initConfig() {
	InitConfig()

	args := map[string]string{
		"username": username,
		"password": password,
		"endpoint": endpoint,
	}

	for key, value := range args {
		if ondemandCmdV.PersistentFlags().Lookup(key).Changed {
			viper.Set(key, value)
		}
		os.Setenv(fmt.Sprintf("VCLOUDAIR_%v", strings.ToUpper(key)), viper.GetString(key))

		if viper.GetString(key) == "" {
			log.Fatalf("missing %v parameter", key)
		}
	}
}

var ondemandplansCmd = &cobra.Command{
	Use:   "plans",
	Short: "plans",
	Long:  `plans`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var ondemandplansgetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   cmdGetPlans,
}

var ondemandservicegroupidsCmd = &cobra.Command{
	Use:   "servicegroupids",
	Short: "servicegroupids",
	Long:  `servicegroupids`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var ondemandservicegroupidsgetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   cmdGetServiceGroupIds,
}

var ondemandinstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "instances",
	Long:  `instances`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var ondemandinstancesgetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   cmdGetInstances,
}

func addCommands() {
	ondemandCmd.AddCommand(ondemandplansCmd)
	ondemandplansCmd.AddCommand(ondemandplansgetCmd)
	ondemandCmd.AddCommand(ondemandservicegroupidsCmd)
	ondemandservicegroupidsCmd.AddCommand(ondemandservicegroupidsgetCmd)
	ondemandCmd.AddCommand(ondemandinstancesCmd)
	ondemandinstancesCmd.AddCommand(ondemandinstancesgetCmd)
}

func authenticate() (client *govcloudair.ODClient, err error) {
	client, err = govcloudair.NewClient()
	if err != nil {
		return client, fmt.Errorf("error with NewClient: %s", err)
	}

	err = client.Authenticate("", "", "", "")
	if err != nil {
		return client, fmt.Errorf("error Authenticating: %s", err)
	}

	return client, err
}

func cmdGetPlans(cmd *cobra.Command, args []string) {
	initConfig()
	client, err := authenticate()
	if err != nil {
		log.Fatal(err)
	}

	planList, err := client.GetPlans()
	if err != nil {
		fmt.Errorf("error Getting plans: %s", err)
	}

	table := table.Table{
		Header:  []string{"Region", "ID", "Name", "ServiceName"},
		Columns: []string{"Region", "ID", "Name", "ServiceName"},
		RowData: reflect.ValueOf(&planList.Plans).Elem(),
	}

	table.PrintTable()
}

func cmdGetServiceGroupIds(cmd *cobra.Command, args []string) {
	initConfig()
	client, err := authenticate()
	if err != nil {
		log.Fatal(err)
	}

	table := table.Table{
		Header:  []string{"ServiceGroupId"},
		RowData: reflect.ValueOf(&client.ServiceGroupIds.ServiceGroupId).Elem(),
	}

	table.PrintColumn()
}

func cmdGetInstances(cmd *cobra.Command, args []string) {
	initConfig()
	client, err := authenticate()
	if err != nil {
		log.Fatal(err)
	}

	instanceList, err := client.GetInstances()
	if err != nil {
		fmt.Errorf("error Getting instances: %s", err)
	}

	table := table.Table{
		Header:  []string{"APIURL", "InstanceAttributes", "Region"},
		Columns: []string{"APIURL", "InstanceAttributes", "Region"},
		RowData: reflect.ValueOf(&instanceList.Instances).Elem(),
	}

	table.PrintTable()
}
