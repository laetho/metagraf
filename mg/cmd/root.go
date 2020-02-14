/*
Copyright 2018-2020 The MetaGraph Authors

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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"os"
)

var MGVersion string
var MGBanner string = "mg " + MGVersion

var (
	All		  bool = false		// Flag for indicating all resources, example delete all resources associated with a application.
	Namespace string
	OName     string // Flag for overriding application name.
	Config    string // Viper config override
	Verbose   bool = false
	Output    bool = false
	Version   string
	Dryrun    bool = false 	// If true do not create
	Branch    string
	Replicas  int = 1			// Flag for adjusting number of replicas.
	BaseEnvs  bool = false
	CVfile    string
	Defaults  bool = false	// Should we hydrate default values in declarative state.
	Format    string
	Template  string			// Command line flag for setting template name
	Suffix    string			// Command line flag for setting mg create ref output file suffix
	Enforce		bool = false	// Boolean flag for articulating enforcement mode instead of inform
	ImageNS		string 			// Image Namespace, used in overriding namespace in container image references
	Registry	string			// Flag for holding a custom container registry
	Tag			string			//
	Context		string			// Flag for setting application context root.
	CreateGlobals bool = false	// Flag for overriding default behaviour of skipping creation of global secrets.
	CVars 		[]string 		// Slice of strings to hold overridden values.
)

// Array of available config keys
var configkeys []string = []string{
	"namespace",
	"user",
	"password",
	"registry",
}

var RootCmd = &cobra.Command{
	Use:   "mg",
	Short: "mg operates on collections of metaGraf's objects.",
	Long: MGBanner + `is a utility that understands the metaGraf
datastructure and help you generate kubernetes primitives`,
	//Run: func(cmd *cobra.Command, args []string) {
	// Do Stuff Here
	//},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&Config, "config", "", "config file (default is $HOME/.config/mg/mg.yaml)")
	RootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
	RootCmd.PersistentFlags().BoolVar(&Output, "output", false, "also output objects")
	RootCmd.PersistentFlags().StringVarP(&Format, "format","o","json", "specify json or yaml, json id default")
	RootCmd.PersistentFlags().StringVar(&Version, "version", "", "Override version in metaGraf specification.")
	RootCmd.PersistentFlags().BoolVar(&Dryrun, "dryrun", false, "do not create objects, only output")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigType("yaml")

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if Config != "" {
		fmt.Printf("Using configfile: %v %v", os.Stdout, Config)
		viper.SetConfigFile(Config)
	} else {
		//fmt.Println(os.Stderr, "Using default config file: ~/.config/mg/config.yaml")
		viper.AddConfigPath(home + "/.config/mg/")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Infof("Failed to read config file: %v", viper.ConfigFileUsed())
	}
}

func Execute() error {
	flag.Parse()
	if err := RootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
