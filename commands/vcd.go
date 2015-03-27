package commands

import (
	"fmt"
	"log"
	"net/url"

	"github.com/emccode/clue"
	"github.com/emccode/govcloudair"
	"github.com/emccode/govcloudair/types/vcav1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmware/govcloudair/types/v56"
	"gopkg.in/yaml.v1"
)

var vcdCmd = &cobra.Command{
	Use:   "vcd",
	Short: "vcd",
	Long:  `vcd services`,
}

var vcdCmdV *cobra.Command

func init() {
	addCommandsVCD()
	vcdCmd.PersistentFlags().StringVar(&username, "username", "", "VCLOUDAIR_USERNAME")
	vcdCmd.PersistentFlags().StringVar(&password, "password", "", "VCLOUDAIR_PASSWORD")
	vcdCmd.PersistentFlags().StringVar(&orgname, "orgname", "", "VCLOUDAIR_ORGNAME")
	vcdCmd.PersistentFlags().StringVar(&insecure, "insecure", "", "VCLOUDAIR_INSECURE")
	vcdCmd.PersistentFlags().StringVar(&vdcname, "vdcname", "", "VCLOUDAIR_VDCNAME")
	vcdCmd.PersistentFlags().StringVar(&vdchref, "vdchref", "", "VCLOUDAIR_VDCHREF")
	vcdCmd.PersistentFlags().StringVar(&sessionuri, "sessionuri", "", "VCLOUDAIR_SESSIONURI")

	vcdCmdV = vcdCmd

	vcdCmd.Run = func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	}

}

var vcdloginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to vCloud Air",
	Long:  `Login to vCloud Air`,
	Run:   cmdLoginVCD,
}

var vcdlogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of vCloud Air",
	Long:  `Logout of vCloud Air and remove temporary token file`,
	Run:   cmdLogout,
}

var vcdvdcCmd = &cobra.Command{
	Use:   "vdc",
	Short: "vdc",
	Long:  `vdc`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var vcdvdcgetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   cmdGetVDC,
}

var vcdvdcuseCmd = &cobra.Command{
	Use:   "use",
	Short: "use",
	Long:  `use`,
	Run:   cmdUseVDC,
}

func addCommandsVCD() {
	vcdCmd.AddCommand(vcdloginCmd)
	vcdCmd.AddCommand(vcdlogoutCmd)
	vcdCmd.AddCommand(vcdvdcCmd)
	vcdvdcCmd.AddCommand(vcdvdcgetCmd)
	vcdvdcCmd.AddCommand(vcdvdcuseCmd)
}

func authenticatevcd(force bool) (client *govcloudair.Client, err error) {
	client, err = govcloudair.NewClient()
	if err != nil {
		return client, fmt.Errorf("error with NewClient: %s", err)
	}

	getValue := clue.GetValue{}
	if err := clue.DecodeGobFile("goair_client", &getValue); err != nil {
		return &govcloudair.Client{}, fmt.Errorf("Problem with client decodeGobFile: %v", err)
	}

	var (
		orgname    = viper.GetString("orgname")
		sessionuri = viper.GetString("sessionuri")
	)

	if force || getValue.VarMap["VCDToken"] == nil {
		instanceAttributes := vcatypes.InstanceAttributes{OrgName: orgname, SessionURI: sessionuri}
		err = client.GetBackendAuthOD(instanceAttributes)
		if err != nil {
			return client, fmt.Errorf("error Authenticating: %s", err)
		}

		err = clue.EncodeGobFile("goair_client", clue.UseValue{
			VarMap: map[string]string{
				"VCDToken":      client.VCDToken,
				"VCDORGHREF":    client.VCDORGHREF.String(),
				"VCDAuthHeader": client.VCDAuthHeader,
			},
		})
	} else {
		client.VCDToken = *getValue.VarMap["VCDToken"]

		orgUri, err := url.ParseRequestURI(*getValue.VarMap["VCDORGHREF"])
		if err != nil {
			return client, fmt.Errorf("cannot parse endpoint coming from VCDORGHREF")
		}

		client.VCDORGHREF = *orgUri
		client.VCDAuthHeader = *getValue.VarMap["VCDAuthHeader"]
	}

	return client, err
}

func cmdLoginVCD(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_vcd", true, map[string]FlagValue{
		"sessionuri": {sessionuri, true, false, ""},
		"orgname":    {orgname, true, false, ""},
	})

	_, err := authenticatevcd(true)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}
	fmt.Println("Successfuly logged in to vCloud Director.")
}

func cmdGetVDC(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_vcd", true, map[string]FlagValue{})
	client, err := authenticatevcd(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	links, err := govcloudair.GetOrgVdc(client, &client.VCDORGHREF)
	if err != nil {
		log.Fatalf("error Getting OrgVdcs: %s", err)
	}

	yamlOutput, err := yaml.Marshal(&links)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))
}

func cmdUseVDC(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_vcd", true, map[string]FlagValue{
		"vdchref": {vdchref, true, false, ""},
		"vdcname": {vdcname, true, false, "vdchref"},
	})
	client, err := authenticatevcd(false)
	if err != nil {
		log.Fatalf("failed authenticating: %s", err)
	}

	links, err := govcloudair.GetOrgVdc(client, &client.VCDORGHREF)

	useLink := types.Link{}
	for _, link := range links {
		if (viper.GetString("vdchref") != "" && link.HREF == viper.GetString("vdchref")) ||
			(viper.GetString("vdcname") != "" && link.Name == viper.GetString("vdcname")) {
			useLink = link
			break
		}
	}

	if useLink.Name != "" {
		err = clue.EncodeGobFile("goair_compute", clue.UseValue{
			VarMap: map[string]string{
				"vdchref":       useLink.HREF,
				"VCDToken":      client.VCDToken,
				"VCDAuthHeader": client.VCDAuthHeader,
				"orghref":       client.VCDORGHREF.String(),
			},
		})

		yamlOutput, err := yaml.Marshal(&useLink)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))

	} else {
		log.Fatalf("Failed to find VDC %v%v", vdcname, vdchref)
	}
}
