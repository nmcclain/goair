package commands

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
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

type UseValue struct {
	VarMap map[string]string
}

type GetValue struct {
	VarMap map[string]*string
}

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

func initConfig(cmd *cobra.Command, flags map[string]FlagValue) {
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

	for key, field := range defaultFlags {
		switch field.persistent {
		case true:
			if cmd.PersistentFlags().Lookup(key).Changed {
				if field.overrideby != "" && (cmd.PersistentFlags().Lookup(field.overrideby).Changed || viper.GetString(field.overrideby) != "") {
					viper.Set(key, "")
					os.Setenv(fmt.Sprintf("VCLOUDAIR_%v", strings.ToUpper(key)), "")
					continue
				} else {
					viper.Set(key, field.value)
				}
			}
		case false:
			if cmd.Flags().Lookup(key).Changed {
				if field.overrideby != "" && (cmd.Flags().Lookup(field.overrideby).Changed || viper.GetString(field.overrideby) != "") {
					viper.Set(key, "")
					os.Setenv(fmt.Sprintf("VCLOUDAIR_%v", strings.ToUpper(key)), "")
					continue
				} else {
					viper.Set(key, field.value)
				}
			}
		default:
		}

		os.Setenv(fmt.Sprintf("VCLOUDAIR_%v", strings.ToUpper(key)), viper.GetString(key))

	}

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

func deleteGobFile(suffix string) (err error) {
	fileLocation := fmt.Sprintf("%vgoair_%v.gob", os.TempDir(), suffix)
	err = os.Remove(fileLocation)
	if err != nil {
		return fmt.Errorf("Problem removing file:", err)
	}
	return nil
}

func encodeGobFile(suffix string, useValue UseValue) (err error) {
	fileLocation := fmt.Sprintf("%vgoair_%v.gob", os.TempDir(), suffix)
	//fmt.Println(fileLocation)
	file, err := os.Create(fileLocation)
	if err != nil {
		return fmt.Errorf("Problem creating file:", err)
	}

	if err = file.Chmod(0600); err != nil {
		return fmt.Errorf("Problem setting persmission onfile:", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("Problem closing file:", err)
		}
	}()

	fileWriter := bufio.NewWriter(file)

	encoder := gob.NewEncoder(fileWriter)
	err = encoder.Encode(useValue)
	//fmt.Println(useValue)
	if err != nil {
		return fmt.Errorf("Problem encoding gob:", err)
	}
	fileWriter.Flush()
	return
}

func decodeGobFile(suffix string, getValue *GetValue) (err error) {
	fileLocation := fmt.Sprintf("%vgoair_%v.gob", os.TempDir(), suffix)
	//fmt.Println(fileLocation)
	file, err := os.Open(fileLocation)
	if err != nil {
		if os.IsExist(err) {
			log.Fatal("Problem opening file:", err)
		} else {
			return nil
		}
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("Problem closing file:", err)
		}
	}()

	fileReader := bufio.NewReader(file)

	decoder := gob.NewDecoder(fileReader)
	err = decoder.Decode(&getValue)
	if err != nil {
		return fmt.Errorf("Problem decoding file:", err)
	}
	return
}

func setGobValues(cmd *cobra.Command, suffix string) (err error) {
	getValue := GetValue{}
	if err := decodeGobFile(suffix, &getValue); err != nil {
		return fmt.Errorf("Problem with decodeGobFile", err)
	}
	for key, _ := range getValue.VarMap {
		lowerKey := strings.ToLower(key)
		if cmd.Flags().Lookup(lowerKey).Changed == false {
			cmd.Flags().Lookup(lowerKey).Value.Set(*getValue.VarMap[key])
			viper.Set(lowerKey, *getValue.VarMap[key])
		}
	}
	return
}
