package commands

//
// import (
// 	"fmt"
// 	"log"
//
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// )
//
// var vappCmd = &cobra.Command{
// 	Use:   "vapp",
// 	Short: "vapp",
// 	Long:  `vapp`,
// }
//
// var vappCmdV *cobra.Command
//
// func init() {
// 	addCommandsVApp()
// 	vappCmd.PersistentFlags().StringVar(&username, "username", "", "VCLOUDAIR_USERNAME")
// 	vappCmd.PersistentFlags().StringVar(&password, "password", "", "VCLOUDAIR_PASSWORD")
// 	vappCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "VCLOUDAIR_ENDPOINT")
// 	viper.SetDefault("endpoint", "https://us-california-1-3.vchs.vmware.com/")
//
// 	vappCmd.Flags().StringVar(&planID, "planid", "", "VCLOUDAIR_PLANID")
// 	vappCmd.Flags().StringVar(&region, "region", "", "VCLOUDAIR_REGION")
// 	vappgetCmd.Flags().StringVar(&planID, "planid", "", "VCLOUDAIR_PLANID")
// 	vappgetCmd.Flags().StringVar(&region, "region", "", "VCLOUDAIR_REGION")
// 	vappgetCmd.Flags().StringVar(&vdcname, "vdcname", "", "VCLOUDAIR_VDCNAME")
// 	vappgetCmd.Flags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")
//
// 	vappCmdV = vappCmd
//
// 	vappCmd.Run = func(cmd *cobra.Command, args []string) {
// 		setGobValues(cmd, "goair_vapp", "")
// 		cmd.Usage()
// 	}
// }
//
//
// func addCommandsVApp() {
// 	vappCmd.AddCommand(vappgetCmd)
// }
//
