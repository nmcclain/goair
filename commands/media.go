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

var mediaCmdV *cobra.Command

func init() {
	addCommandsMedia()
	mediaCmd.Flags().StringVar(&medianame, "medianame", "", "VCLOUDAIR_MEDIANAME")
	mediagetCmd.Flags().StringVar(&medianame, "medianame", "", "VCLOUDAIR_MEDIANAME")

	mediaCmdV = mediaCmd

	mediaCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goair_compute", "")
		cmd.Usage()
	}
}

func addCommandsMedia() {
	mediaCmd.AddCommand(mediagetCmd)
}

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "media",
	Long:  `media`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var mediagetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a media",
	Long:  `Get a media`,
	Run:   cmdGetMedia,
}

func cmdGetMedia(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"vdchref": {vdchref, true, false, ""},
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
	err = vdc.Refresh()
	if err != nil {
		log.Fatalf("err refreshing vdc: %v", err)
	}

	media, err := vdc.GetMedia()
	if err != nil {
		log.Fatalf("Error problem getting media from vdc: %v", err)
	}

	yamlOutput, err := yaml.Marshal(&media)
	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))
	return

}

// resourceEntities := types.ResourceEntities{}
// for _, link := range links {
//   if (viper.GetString("vdchref") != "" && link.HREF == viper.GetString("vdchref")) ||
//     (viper.GetString("vdcname") != "" && link.Name == viper.GetString("vdcname")) {
//     useLink = link
//     break
//   }
// }
//
