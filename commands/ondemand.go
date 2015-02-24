package commands

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/emccode/clue"
	"github.com/emccode/gotablethis"
	"github.com/emccode/govcloudair"
	"github.com/emccode/govcloudair/types/vcav1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

var ondemandCmd = &cobra.Command{
	Use:   "ondemand",
	Short: "ondemand",
	Long:  `ondemand services`,
}

var ondemandCmdV *cobra.Command

func init() {
	addCommandsOnDemand()
	ondemandCmd.PersistentFlags().StringVar(&username, "username", "", "VCLOUDAIR_USERNAME")
	ondemandCmd.PersistentFlags().StringVar(&password, "password", "", "VCLOUDAIR_PASSWORD")
	ondemandCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "VCLOUDAIR_ENDPOINT")
	viper.SetDefault("endpoint", "https://us-california-1-3.vchs.vmware.com/api")

	ondemandbillablecostsCmd.Flags().StringVar(&serviceGroupID, "servicegroupid", "", "VCLOUDAIR_SERVICEGROUPID")
	ondemandbillablecostsgetCmd.Flags().StringVar(&serviceGroupID, "servicegroupid", "", "VCLOUDAIR_SERVICEGROUPID")

	ondemandCmdV = ondemandCmd

	ondemandCmd.Run = func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	}

}

var ondemandloginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to vCloud Air",
	Long:  `Login to vCloud Air`,
	Run:   cmdLogin,
}

var ondemandlogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of vCloud Air",
	Long:  `Logout of vCloud Air and remove temporary token file`,
	Run:   cmdLogout,
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

var ondemandinstancesnewCmd = &cobra.Command{
	Use:   "new",
	Short: "new",
	Long:  `new`,
	Run:   cmdNewInstance,
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
		ondemandCmd.Flags().StringVar(&serviceGroupID, "servicegroupid", "", "VCLOUDAIR_SERVICEGROUPID")
		cmd.Usage()
	},
}

var ondemandbillablecostsCmd = &cobra.Command{
	Use:   "costs",
	Short: "costs",
	Long:  `costs`,
	Run: func(cmd *cobra.Command, args []string) {
		ondemandCmd.Flags().StringVar(&serviceGroupID, "servicegroupid", "", "VCLOUDAIR_SERVICEGROUPID")
		cmd.Usage()
	},
}

var ondemandbillablecostsgetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   cmdGetBillableCosts,
}

func addCommandsOnDemand() {
	ondemandCmd.AddCommand(ondemandloginCmd)
	ondemandCmd.AddCommand(ondemandlogoutCmd)
	ondemandCmd.AddCommand(ondemandplansCmd)
	ondemandplansCmd.AddCommand(ondemandplansgetCmd)
	ondemandCmd.AddCommand(ondemandservicegroupidsCmd)
	ondemandservicegroupidsCmd.AddCommand(ondemandservicegroupidsgetCmd)
	ondemandCmd.AddCommand(ondemandinstancesCmd)
	ondemandinstancesCmd.AddCommand(ondemandinstancesgetCmd)
	// ondemandinstancesCmd.AddCommand(ondemandinstancesnewCmd)
	ondemandCmd.AddCommand(ondemandusersCmd)
	ondemandusersCmd.AddCommand(ondemandusersgetCmd)
	ondemandCmd.AddCommand(ondemandbillableCmd)
	ondemandbillableCmd.AddCommand(ondemandbillablecostsCmd)
	ondemandbillablecostsCmd.AddCommand(ondemandbillablecostsgetCmd)
}

func authenticate(force bool) (client *govcloudair.Client, err error) {
	client, err = govcloudair.NewClient()
	if err != nil {
		return client, fmt.Errorf("error with NewClient: %s", err)
	}

	getValue := clue.GetValue{}
	if err := clue.DecodeGobFile("goair_client", &getValue); err != nil {
		return &govcloudair.Client{}, fmt.Errorf("Problem with client decodeGobFile: %v", err)
	}

	if force || getValue.VarMap["VAToken"] == nil {
		err = client.Authenticate("", "", "", "")
		if err != nil {
			return client, fmt.Errorf("error Authenticating: %s", err)
		}

		err = clue.EncodeGobFile("goair_client", clue.UseValue{
			VarMap: map[string]string{
				"VAToken": client.VAToken,
			},
		})
	} else {
		fmt.Println(client.VAToken)
		client.VAToken = *getValue.VarMap["VAToken"]
	}

	return client, err
}

func cmdLogin(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_ondemand", true, map[string]FlagValue{})
	_, err := authenticate(true)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}
	fmt.Println("Successfuly logged in to vCloud Air On Demand.")
}

func cmdLogout(cmd *cobra.Command, args []string) {
	err := clue.DeleteGobFile("goair_client")
	if err != nil {
		if os.IsExist(err) {
			log.Fatalf("failed to delete client gob file: %s", err)
		}
	}
}

func cmdGetPlans(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_ondemand", true, map[string]FlagValue{})
	client, err := authenticate(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	planList, err := client.GetPlans()
	if err != nil {
		log.Fatalf("error Getting plans: %s", err)
	}

	// table := table.Table{
	// 	Header:  []string{"Region", "ID", "Name", "ServiceName"},
	// 	Columns: []string{"Region", "ID", "Name", "ServiceName"},
	// 	RowData: reflect.ValueOf(&planList.Plans).Elem(),
	// }
	//
	// table.PrintTable()

	yamlOutput, err := yaml.Marshal(&planList)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))
}

func cmdGetServiceGroupIds(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_ondemand", true, map[string]FlagValue{})
	client, err := authenticate(true)
	if err != nil {
		log.Fatal(err)
	}

	table := gotablethis.Table{
		Header:  []string{"ServiceGroupID"},
		RowData: reflect.ValueOf(&client.ServiceGroupIds.ServiceGroupID).Elem(),
	}
	table.PrintColumn()
}

func cmdGetInstances(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_ondemand", true, map[string]FlagValue{})
	client, err := authenticate(false)
	if err != nil {
		log.Fatal(err)
	}

	instanceList, err := client.GetInstances()
	if err != nil {
		log.Fatalf("error Getting instances: %s", err)
	}

	// for _, arg := range instanceList.Instances {
	// 	table := gotablethis.Table{
	// 		RowData: reflect.ValueOf(&arg).Elem(),
	// 	}
	// 	table.PrintKeyValueTable()
	// 	fmt.Println()
	// }
	yamlOutput, err := yaml.Marshal(&instanceList.Instances)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))

}

func cmdNewInstance(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_ondemand", true, map[string]FlagValue{})
	client, err := authenticate(false)
	if err != nil {
		log.Fatal(err)
	}

	instanceSpecParams := vcatypes.InstanceSpecParams{
		Name:           "testing",
		PlanID:         "41400e74-4445-49ef-90a4-98da4ccfb16c",
		ServiceGroupID: "4fde19a4-7621-428e-b190-dd4db2e158cd",
	}

	instance, err := client.NewInstance(instanceSpecParams)
	if err != nil {
		log.Fatalf("error Getting instances: %s", err)
	}

	table := gotablethis.Table{
		RowData: reflect.ValueOf(&instance).Elem(),
	}
	table.PrintKeyValueTable()
	fmt.Println()

}

func cmdGetUsers(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_ondemand", true, map[string]FlagValue{})
	client, err := authenticate(false)
	if err != nil {
		log.Fatal(err)
	}

	users, err := client.GetUsers()
	if err != nil {
		log.Fatalf("error Getting users: %s", err)
	}

	for _, arg := range users.User {
		table := gotablethis.Table{
			RowData: reflect.ValueOf(&arg).Elem(),
		}
		table.PrintKeyValueTable()
		fmt.Println()
	}

}

func cmdGetBillableCosts(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_ondemand", true, map[string]FlagValue{
		"servicegroupid": {serviceGroupID, true, false, ""},
	})
	client, err := authenticate(false)
	if err != nil {
		log.Fatal(err)
	}

	billableCosts, err := client.GetBillableCosts(serviceGroupID)
	if err != nil {
		log.Fatalf("error Getting billable costs: %s", err)

	}

	for _, arg := range billableCosts.Cost {
		table := gotablethis.Table{
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

	table := gotablethis.Table{
		RowData: reflect.ValueOf(&leftOver).Elem(),
	}
	table.PrintKeyValueTable()
	fmt.Println()

}
