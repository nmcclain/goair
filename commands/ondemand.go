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
	Run:   GetPlans,
}

func addCommands() {
	ondemandCmd.AddCommand(ondemandplansCmd)
	ondemandplansCmd.AddCommand(ondemandplansgetCmd)
}

func GetPlans(cmd *cobra.Command, args []string) {
	initConfig()
	client, err := govcloudair.NewClient()
	if err != nil {
		fmt.Errorf("error with NewClient: %s", err)
	}

	err = client.Authenticate("", "", "", "")
	if err != nil {
		fmt.Errorf("error Authenticating: %s", err)
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

	table.Printtable()
}
