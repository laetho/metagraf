/*
Copyright 2018 The MetaGraph Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

const Banner string = "mg (metaGraf) - "


var (
	Namespace string
)

// Array of available config keys
var configkeys []string = []string{
	"namespace",
	"user",
	"password",
	"registry",
}

// Flags
var Config	string			// Viper config override
var Verbose	bool = false
var Output	bool = false
var Version	string
var Dryrun	bool = false // If true do not create
var Branch	string

var RootCmd = &cobra.Command{
	Use:   "mg",
	Short: "mg operates on collections of metaGraf's objects.",
	Long: Banner + `is a utility that understands the metaGraf
datastructure and help you generate kubernetes primitives`,
	//Run: func(cmd *cobra.Command, args []string) {
	// Do Stuff Here
	//},
}

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	RootCmd.PersistentFlags().StringVar(&Config, "config", "", "config file (default is $HOME/.config/mg/mg.yaml)")
	RootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
	RootCmd.PersistentFlags().BoolVar(&Output, "output", false, "also output objects in json")
	RootCmd.PersistentFlags().StringVar(&Version, "version", "", "Override version in metaGraf specification.")
	RootCmd.PersistentFlags().BoolVar(&Dryrun, "dryrun", false, "do not create objects, only output")

}

func initConfig() {
	viper.SetConfigType("yaml")

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(home + "/.config/mg/")
	viper.SetConfigName("config")

	if len(Config) > 0 {
		viper.SetConfigFile(Config)
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		glog.Infof("Failed to read config file:", viper.ConfigFileUsed())
	}
}

func Execute() error {
	initConfig()

	if err := RootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
