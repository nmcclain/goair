package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emccode/clue"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	CfgFile        string
	username       string
	password       string
	endpoint       string
	serviceGroupId string
	planID         string
	region         string
)

type FlagValue struct {
	value      string
	mandatory  bool
	persistent bool
	overrideby string
}

var GoairCmd = &cobra.Command{
	Use: "goair",
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig()
		cmd.Usage()
	},
}

func Exec() {
	AddCommands()
	GoairCmd.Execute()
}

func AddCommands() {
	GoairCmd.AddCommand(ondemandCmd)
	GoairCmd.AddCommand(computeCmd)
}

var goairCmdV *cobra.Command

func init() {
	GoairCmd.PersistentFlags().StringVar(&CfgFile, "Config", "", "config file (default is $HOME/goair/config.yaml)")
	goairCmdV = GoairCmd
}

func initConfig(cmd *cobra.Command, suffix string, checkValues bool, flags map[string]FlagValue) {
	InitConfig()

	defaultFlags := map[string]FlagValue{
		"username": {username, true, false, ""},
		"password": {password, true, false, ""},
		"endpoint": {endpoint, true, false, ""},
	}

	for key, field := range flags {
		defaultFlags[key] = field
	}

	fieldsMissing := make([]string, 0)
	fieldsMissingRemove := make([]string, 0)

	cmdFlags := &pflag.FlagSet{}

	for key, field := range defaultFlags {
		viper.BindEnv(key)

		switch field.persistent {
		case true:
			cmdFlags = cmd.PersistentFlags()
		case false:
			cmdFlags = cmd.Flags()
		default:
		}

		if cmdFlags.Lookup(key).Changed {
			if field.overrideby != "" {
				if cmdFlags.Lookup(field.overrideby).Changed {
					viper.Set(key, "")
					continue
				}
			}
			viper.Set(key, cmdFlags.Lookup(key).Value)
		} else {
			if field.overrideby != "" && cmdFlags.Lookup(field.overrideby).Changed == false && viper.GetString(field.overrideby) == "" {
				if viper.GetString(key) == "" {
					if err := setGobValues(cmd, "goair_compute", key); err != nil {
						log.Fatal(err)
					}
				}
			} else {
				if field.overrideby == "" {
					if viper.GetString(key) == "" {
						if err := setGobValues(cmd, "goair_compute", key); err != nil {
							log.Fatal(err)
						}
						for removeKey, field := range defaultFlags {
							if key == field.overrideby {
								viper.Set(removeKey, "")
							}
						}
					}
				}
			}
		}
	}

	if checkValues {
		for key, field := range defaultFlags {
			if field.mandatory == true {
				if viper.GetString(key) != "" && (field.overrideby != "" && viper.GetString(field.overrideby) == "") {
					fieldsMissingRemove = append(fieldsMissingRemove, field.overrideby)
				} else {
					if viper.GetString(key) == "" && (field.overrideby != "" && viper.GetString(field.overrideby) == "") {
						fieldsMissing = append(fieldsMissing, key)
					}
				}
			}
		}

	Loop1:
		for _, fieldMissingRemove := range fieldsMissingRemove {
			for i, fieldMissing := range fieldsMissing {
				if fieldMissing == fieldMissingRemove {
					fieldsMissing = append(fieldsMissing[:i], fieldsMissing[i+1:]...)
					break Loop1
				}
			}
		}

		if len(fieldsMissing) != 0 {
			log.Fatalf("missing parameter: %v", strings.Join(fieldsMissing, ", "))
		}
	}

	for key, _ := range defaultFlags {
		if viper.GetString(key) != "" {
			os.Setenv(fmt.Sprintf("VCLOUDAIR_%v", strings.ToUpper(key)), viper.GetString(key))
		}
		//fmt.Println(viper.GetString(key))
	}

}

func InitConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/goair")
	viper.AddConfigPath("$HOME/.goair/")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VCLOUDAIR")
}

func setGobValues(cmd *cobra.Command, suffix string, field string) (err error) {
	getValue := clue.GetValue{}
	if err := clue.DecodeGobFile(suffix, &getValue); err != nil {
		return fmt.Errorf("Problem with decodeGobFile", err)
	}
	for key, _ := range getValue.VarMap {
		lowerKey := strings.ToLower(key)
		if field != "" && field != lowerKey {
			continue
		}
		viper.Set(lowerKey, *getValue.VarMap[key])
	}
	return
}
