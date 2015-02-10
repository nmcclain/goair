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

var ondemandCmd = &cobra.Command{
	Use:   "ondemand",
	Short: "ondemand",
	Long:  `Get plans`,
}

var ondemandCmdV *cobra.Command

var (
	username       string
	password       string
	endpoint       string
	serviceGroupId string
)

func init() {
	addCommands()
	ondemandCmd.PersistentFlags().StringVar(&username, "username", "", "VCLOUDAIR_USERNAME")
	ondemandCmd.PersistentFlags().StringVar(&password, "password", "", "VCLOUDAIR_PASSWORD")
	ondemandCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "VCLOUDAIR_ENDPOINT")
	viper.SetDefault("endpoint", "https://us-california-1-3.vchs.vmware.com/")

	ondemandbillablecostsCmd.Flags().StringVar(&serviceGroupId, "servicegroupid", "", "VCLOUDAIR_SERVICEGROUPID")
	ondemandbillablecostsgetCmd.Flags().StringVar(&serviceGroupId, "servicegroupid", "", "VCLOUDAIR_SERVICEGROUPID")

	ondemandCmdV = ondemandCmd

	ondemandCmd.Run = func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	}

}

func initConfig(cmd *cobra.Command, flags map[string]FlagValue) {
	InitConfig()

	defaultFlags := map[string]FlagValue{
		"username": {username, true, false},
		"password": {password, true, false},
		"endpoint": {endpoint, true, false},
	}

	for key, field := range flags {
		defaultFlags[key] = field
	}

	for key, field := range defaultFlags {
		switch field.persistent {
		case true:
			if cmd.PersistentFlags().Lookup(key).Changed {
				viper.Set(key, field.value)
			}
		case false:
			if cmd.Flags().Lookup(key).Changed {
				viper.Set(key, field.value)
			}
		default:
		}

		os.Setenv(fmt.Sprintf("VCLOUDAIR_%v", strings.ToUpper(key)), viper.GetString(key))

		if viper.GetString(key) == "" && field.mandatory == true {
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

var ondemandusersCmd = &cobra.Command{
	Use:   "users",
	Short: "users",
	Long:  `users`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var ondemandusersgetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   cmdGetUsers,
}

var ondemandbillableCmd = &cobra.Command{
	Use:   "billable",
	Short: "billable",
	Long:  `billable`,
	Run: func(cmd *cobra.Command, args []string) {
		ondemandCmd.Flags().StringVar(&serviceGroupId, "servicegroupid", "", "VCLOUDAIR_SERVICEGROUPID")
		cmd.Usage()
	},
}

var ondemandbillablecostsCmd = &cobra.Command{
	Use:   "costs",
	Short: "costs",
	Long:  `costs`,
	Run: func(cmd *cobra.Command, args []string) {
		ondemandCmd.Flags().StringVar(&serviceGroupId, "servicegroupid", "", "VCLOUDAIR_SERVICEGROUPID")
		cmd.Usage()
	},
}

var ondemandbillablecostsgetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   cmdGetBillableCosts,
}

func addCommands() {
	ondemandCmd.AddCommand(ondemandplansCmd)
	ondemandplansCmd.AddCommand(ondemandplansgetCmd)
	ondemandCmd.AddCommand(ondemandservicegroupidsCmd)
	ondemandservicegroupidsCmd.AddCommand(ondemandservicegroupidsgetCmd)
	ondemandCmd.AddCommand(ondemandinstancesCmd)
	ondemandinstancesCmd.AddCommand(ondemandinstancesgetCmd)
	ondemandCmd.AddCommand(ondemandusersCmd)
	ondemandusersCmd.AddCommand(ondemandusersgetCmd)
	ondemandCmd.AddCommand(ondemandbillableCmd)
	ondemandbillableCmd.AddCommand(ondemandbillablecostsCmd)
	ondemandbillablecostsCmd.AddCommand(ondemandbillablecostsgetCmd)
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
	initConfig(cmd, map[string]FlagValue{})

	client, err := authenticate()
	if err != nil {
		log.Fatal(err)
	}

	planList, err := client.GetPlans()
	if err != nil {
		log.Fatalf("error Getting plans: %s", err)
	}

	table := table.Table{
		Header:  []string{"Region", "ID", "Name", "ServiceName"},
		Columns: []string{"Region", "ID", "Name", "ServiceName"},
		RowData: reflect.ValueOf(&planList.Plans).Elem(),
	}

	table.PrintTable()
}

func cmdGetServiceGroupIds(cmd *cobra.Command, args []string) {
	initConfig(cmd, map[string]FlagValue{})
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
	initConfig(cmd, map[string]FlagValue{})
	client, err := authenticate()
	if err != nil {
		log.Fatal(err)
	}

	instanceList, err := client.GetInstances()
	if err != nil {
		log.Fatalf("error Getting instances: %s", err)
	}

	for _, arg := range instanceList.Instances {
		table := table.Table{
			RowData: reflect.ValueOf(&arg).Elem(),
		}
		table.PrintKeyValueTable()
		fmt.Println()
	}

}

func cmdGetUsers(cmd *cobra.Command, args []string) {
	initConfig(cmd, map[string]FlagValue{})
	client, err := authenticate()
	if err != nil {
		log.Fatal(err)
	}

	users, err := client.GetUsers()
	if err != nil {
		log.Fatalf("error Getting users: %s", err)
	}

	for _, arg := range users.User {
		table := table.Table{
			RowData: reflect.ValueOf(&arg).Elem(),
		}
		table.PrintKeyValueTable()
		fmt.Println()
	}

}

func cmdGetBillableCosts(cmd *cobra.Command, args []string) {
	initConfig(cmd, map[string]FlagValue{
		"servicegroupid": {serviceGroupId, true, false},
	})
	client, err := authenticate()
	if err != nil {
		log.Fatal(err)
	}

	billableCosts, err := client.GetBillableCosts(serviceGroupId)
	if err != nil {
		log.Fatalf("error Getting billable costs: %s", err)

	}

	for _, arg := range billableCosts.Cost {
		table := table.Table{
			RowData: reflect.ValueOf(&arg).Elem(),
		}
		table.PrintKeyValueTable()
		fmt.Println()
	}

	type tempst struct {
		Currency       string
		LastUpdateTime string
	}
	leftOver := tempst{
		Currency:       billableCosts.Currency,
		LastUpdateTime: billableCosts.LastUpdateTime,
	}

	table := table.Table{
		RowData: reflect.ValueOf(&leftOver).Elem(),
	}
	table.PrintKeyValueTable()
	fmt.Println()

}
