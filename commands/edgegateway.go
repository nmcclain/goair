package commands

import (
	"fmt"
	"log"
	"net/url"

	"github.com/emccode/govcloudair"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	types "github.com/vmware/govcloudair/types/v56"
	"gopkg.in/yaml.v1"
)

//PublicIPInfo holds the NAT mapping information
type PublicIPInfo struct {
	Type           string
	SNATRuleExists bool
	DNATRuleExists bool
	ExternalIP     string
	InternalIP     string
}

var edgegatewayCmdV *cobra.Command

func init() {
	addCommandsEdgeGateway()
	edgegatewaynewnatCmd.Flags().StringVar(&externalip, "externalip", "", "VCLOUDAIR_EXTERNALIP")
	edgegatewaynewnatCmd.Flags().StringVar(&internalip, "internalip", "", "VCLOUDAIR_INTERNALIP")
	edgegatewaynewnatCmd.Flags().StringVar(&description, "description", "", "VCLOUDAIR_DESCRIPTION")
	edgegatewaynewnatCmd.Flags().StringVar(&runasync, "runasync", "", "VCLOUDAIR_RUNASYNC")
	edgegatewayremovenatCmd.Flags().StringVar(&externalip, "externalip", "", "VCLOUDAIR_EXTERNALIP")
	edgegatewayremovenatCmd.Flags().StringVar(&internalip, "internalip", "", "VCLOUDAIR_INTERNALIP")
	edgegatewayremovenatCmd.Flags().StringVar(&runasync, "runasync", "", "VCLOUDAIR_RUNASYNC")
	edgegatewaynewfirewallCmd.Flags().StringVar(&sourceip, "sourceip", "", "VCLOUDAIR_SOURCEIP")
	edgegatewaynewfirewallCmd.Flags().StringVar(&sourceport, "sourceport", "", "VCLOUDAIR_SOURCEPORT")
	edgegatewaynewfirewallCmd.Flags().StringVar(&destinationip, "destinationip", "", "VCLOUDAIR_DESTINATIONIP")
	edgegatewaynewfirewallCmd.Flags().StringVar(&destinationport, "destinationport", "", "VCLOUDAIR_DESTINATIONPORT")
	edgegatewaynewfirewallCmd.Flags().StringVar(&description, "description", "", "VCLOUDAIR_DESCRIPTION")

	edgegatewayCmdV = edgegatewayCmd

	edgegatewayCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goair_compute", "")
		cmd.Usage()
	}
}

func addCommandsEdgeGateway() {
	edgegatewayCmd.AddCommand(edgegatewaygetCmd)
	edgegatewayCmd.AddCommand(edgegatewaynewnatCmd)
	edgegatewayCmd.AddCommand(edgegatewayremovenatCmd)
	edgegatewayCmd.AddCommand(edgegatewaynewfirewallCmd)
}

var edgegatewayCmd = &cobra.Command{
	Use:   "edgegateway",
	Short: "edgegateway",
	Long:  `edgegateway`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var edgegatewaygetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get edgegateway",
	Long:  `Get edgegateway`,
	Run:   cmdGetEdgeGateway,
}

var edgegatewaynewnatCmd = &cobra.Command{
	Use:   "new-natrule",
	Short: "Create NAT statement on edgegateway",
	Long:  `Create NAT statement on edgegateway`,
	Run:   cmdNewNatEdgeGateway,
}

var edgegatewayremovenatCmd = &cobra.Command{
	Use:   "remove-natrule",
	Short: "Remove NAT statement on edgegateway",
	Long:  `Remove NAT statement on edgegateway`,
	Run:   cmdRemoveNatEdgeGateway,
}

var edgegatewaynewfirewallCmd = &cobra.Command{
	Use:   "new-firewallrule",
	Short: "Create firewall rule on edgegateway",
	Long:  `Create firewall rule on edgegateway`,
	Run:   cmdNewFirewallEdgeGateway,
}

func cmdGetEdgeGateway(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
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
	err = vdc.Refresh()
	if err != nil {
		log.Fatalf("err refreshing vdc: %v", err)
	}

	edgeGateway, err := vdc.FindEdgeGateway("")
	if err != nil {
		log.Fatalf("err getting edge gateway: %v", err)
	}

	if len(args) == 0 {
		yamlOutput, err := yaml.Marshal(&edgeGateway)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	var yamlOutput []byte
	switch args[0] {
	case "natrule":
		natRules := edgeGateway.EdgeGateway.Configuration.EdgeGatewayServiceConfiguration.NatService.NatRule
		yamlOutput, err = yaml.Marshal(&natRules)
	case "gatewayinterface":
		gatewayInterfaces := edgeGateway.EdgeGateway.Configuration.GatewayInterfaces.GatewayInterface
		yamlOutput, err = yaml.Marshal(&gatewayInterfaces)
	case "publicip":
		gatewayInterfaces := edgeGateway.EdgeGateway.Configuration.GatewayInterfaces.GatewayInterface
		ipRanges := types.IPRanges{}
		for _, gatewayInterface := range gatewayInterfaces {
			if gatewayInterface.InterfaceType == "uplink" {
				if gatewayInterface.SubnetParticipation.IPRanges != nil {
					for _, ipRange := range gatewayInterface.SubnetParticipation.IPRanges.IPRange {
						ipRanges.IPRange = append(ipRanges.IPRange, ipRange)
					}
				}
			}
		}

		var publicIPInfoMap map[string]*PublicIPInfo
		publicIPInfoMap = make(map[string]*PublicIPInfo)

		for _, ipRange := range ipRanges.IPRange {
			publicIPInfoMap[ipRange.StartAddress] = &PublicIPInfo{"", false, false, ipRange.StartAddress, ""}
		}

		natRules := edgeGateway.EdgeGateway.Configuration.EdgeGatewayServiceConfiguration.NatService.NatRule
		for _, natRule := range natRules {
			switch natRule.RuleType {
			case "DNAT":
				if publicIPInfoMap[natRule.GatewayNatRule.OriginalIP] != nil {
					publicIPInfoMap[natRule.GatewayNatRule.OriginalIP].DNATRuleExists = true
				}
			case "SNAT":
				if publicIPInfoMap[natRule.GatewayNatRule.TranslatedIP] != nil {
					publicIPInfoMap[natRule.GatewayNatRule.TranslatedIP].SNATRuleExists = true
					publicIPInfoMap[natRule.GatewayNatRule.TranslatedIP].InternalIP = natRule.GatewayNatRule.OriginalIP
				}
			}
		}

		publicIPInfo := make([]PublicIPInfo, 0, len(publicIPInfoMap))

		for _, value := range publicIPInfoMap {
			if value.ExternalIP != "" && value.InternalIP != "" {
				value.Type = "1to1"
			}
			publicIPInfo = append(publicIPInfo, *value)
		}

		yamlOutput, err = yaml.Marshal(&publicIPInfo)

	case "firewall":
		firewallService := edgeGateway.EdgeGateway.Configuration.EdgeGatewayServiceConfiguration.FirewallService
		yamlOutput, err = yaml.Marshal(&firewallService)

	default:
		log.Fatalf("need to specify proper parameter after get natrule|gatewayinterface|publicip")
	}

	if err != nil {
		log.Fatalf("error marshaling: %s", err)
	}
	fmt.Println(string(yamlOutput))
	return

}

func cmdNewNatEdgeGateway(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
		"externalip":         {externalip, true, false, ""},
		"internalip":         {internalip, true, false, ""},
		"description":        {description, false, false, ""},
		"runasync":           {runasync, false, false, ""},
		"instanceAttributes": {instanceAttributes, true, false, ""},
	})

	if args[0] != "1to1" {
		log.Fatalf("Missing type of NAT, currently 1to1 is supported.")
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
	err = vdc.Refresh()
	if err != nil {
		log.Fatalf("err refreshing vdc: %v", err)
	}

	edgeGateway, err := vdc.FindEdgeGateway("")
	if err != nil {
		log.Fatalf("err getting edge gateway: %v", err)
	}

	task, err := edgeGateway.Create1to1Mapping(internalip, externalip, description, false, false)
	if err != nil {
		log.Fatalf("err creating 1 to 1 mapping: %v", err)
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

func cmdRemoveNatEdgeGateway(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
		"externalip":         {externalip, true, false, ""},
		"internalip":         {internalip, true, false, ""},
		"runasync":           {runasync, false, false, ""},
		"instanceAttributes": {instanceAttributes, true, false, ""},
	})

	if args[0] != "1to1" {
		log.Fatalf("Missing type of NAT, currently 1to1 is supported.")
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
	err = vdc.Refresh()
	if err != nil {
		log.Fatalf("err refreshing vdc: %v", err)
	}

	edgeGateway, err := vdc.FindEdgeGateway("")
	if err != nil {
		log.Fatalf("err getting edge gateway: %v", err)
	}

	task, err := edgeGateway.Remove1to1Mapping(internalip, externalip)
	if err != nil {
		log.Fatalf("err creating 1 to 1 mapping: %v", err)
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

func cmdNewFirewallEdgeGateway(cmd *cobra.Command, args []string) {
	initConfig(cmd, "goair_compute", true, map[string]FlagValue{
		"planid":             {planID, true, false, ""},
		"region":             {region, true, false, "planid"},
		"vdchref":            {vdchref, true, false, ""},
		"sourceip":           {sourceip, true, false, ""},
		"sourceport":         {sourceport, true, false, ""},
		"destinationip":      {destinationip, true, false, ""},
		"destinationport":    {destinationport, true, false, ""},
		"description":        {description, true, false, ""},
		"runasync":           {runasync, false, false, ""},
		"instanceAttributes": {instanceAttributes, true, false, ""},
	})

	if len(args) > 0 && args[0] != "1to1" {
		log.Fatalf("Missing type of NAT, currently 1to1 is supported.")
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
	err = vdc.Refresh()
	if err != nil {
		log.Fatalf("err refreshing vdc: %v", err)
	}

	edgeGateway, err := vdc.FindEdgeGateway("")
	if err != nil {
		log.Fatalf("err getting edge gateway: %v", err)
	}

	//sourceip sourceport destinationip destinationport
	fwin := &types.FirewallRule{
		Description: description,
		IsEnabled:   true,
		Policy:      "allow",
		Protocols: &types.FirewallRuleProtocols{
			Any: true,
		},
		DestinationPortRange: "Any",
		DestinationIP:        destinationip,
		SourcePortRange:      "Any",
		SourceIP:             "Any",
		EnableLogging:        false,
	}

	newedgeconfig := edgeGateway.EdgeGateway.Configuration.EdgeGatewayServiceConfiguration

	newedgeconfig.FirewallService.FirewallRule = append(newedgeconfig.FirewallService.FirewallRule, fwin)

	task, err := edgeGateway.UpdateFirewall(newedgeconfig)
	if err != nil {
		log.Fatalf("error updating firewall: %v", err)
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
