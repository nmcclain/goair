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

func encodeGobFile(suffix string, useValue UseValue) {
	fileLocation := fmt.Sprintf("%vgoair_%v.gob", os.TempDir(), suffix)
	fmt.Println(fileLocation)
	file, err := os.Create(fileLocation)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fileWriter := bufio.NewWriter(file)

	encoder := gob.NewEncoder(fileWriter)
	err = encoder.Encode(useValue)
	fmt.Println(useValue)
	if err != nil {
		log.Fatal(err)
	}
	fileWriter.Flush()
}

func decodeGobFile(suffix string, getValue *GetValue) {
	fileLocation := fmt.Sprintf("%vgoair_%v.gob", os.TempDir(), suffix)
	fmt.Println(fileLocation)
	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fileReader := bufio.NewReader(file)

	decoder := gob.NewDecoder(fileReader)
	err = decoder.Decode(&getValue)
	if err != nil {
		log.Fatal("decode error:", err)
	}
}
