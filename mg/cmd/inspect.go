package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"os"
)

func init() {
	RootCmd.AddCommand(InspectCmd)
	InspectCmd.Flags().BoolVar(&Enforce, "enforce",false, "Enforce findings, defaults to false and informs only.")
	InspectCmd.AddCommand(InspectPropertiesCmd)
	InspectPropertiesCmd.Flags().StringVar(&CVfile, "cvfile","", "File with component configuration values. (key=value pairs)")
}

var InspectCmd = &cobra.Command{
	Use:   "inspect <metaGraf>",
	Short: "inspect a metaGraf specification",
	Long:  `inspect a metaGraf specification and list objects that will be created or patched.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Active project is:", viper.Get("namespace"))
			fmt.Println("Missing path to metaGraf specification")
			return
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				fmt.Println("Namespace must be supplied")
				os.Exit(1)
			}
		}

		mg := metagraf.Parse(args[0])
		modules.Variables = mg.GetProperties()
		OverrideProperties(modules.Variables)
		log.V(2).Info("Current MGProperties: ", modules.Variables)

		name := modules.Name(&mg)
		for k, v := range modules.Variables {
			fmt.Println(name, "Variable:", k, v)
		}

		modules.InspectSecrets(&mg)
		modules.InspectConfigMaps(&mg)
	},
}

var InspectPropertiesCmd = &cobra.Command{
	Use:   "properties <metaGraf>",
	Short: "inspect a metaGraf specification against a properties file",
	Long:  `inspect a metaGraf specification against a properties file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Infof("Active project is:", viper.Get("namespace"))
			log.Errorf("Missing path to metaGraf specification")
			os.Exit(1)
		}

		if len(args) < 2 {
			log.Errorf("Missing path to properties file")
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				fmt.Println("Namespace must be supplied")
				os.Exit(1)
			}
		}
		mg := metagraf.Parse(args[0])
		CVfile = args[1]

		mgprops := mg.GetProperties()
		fileprops := PropertiesFromFile(mgprops)
		confvars := fileprops.SourceKeyMap(false)
		reqvars := mgprops.GetRequired().SourceKeyMap(true)

		log.V(1).Info("Addressable Variables:", mg.GetProperties())
		log.V(1).Info("Required Variables", reqvars)
		log.V(1).Info("Config Variables: ", confvars)


		// 1. Find required vars not in variables form the properties file.
		// 2. Find configvars not in specification.

		fail := false
		for key,_ := range reqvars {
			if _, ok := fileprops[key]; !ok {
				fail = true
				fmt.Printf("Required key: %v, is missing from %v\n", key, CVfile)
			}
		}

		for key,_ := range confvars {
			if _, ok := mgprops[key]; !ok {
				fail = true
				fmt.Printf("%v is an invalid configuration key for this metaGraf specification.\n", key)
			}
		}
		if fail {
			os.Exit(1)
		}
		fmt.Printf("The %v configuration is valid for this metaGraf specification.\n", CVfile)
		os.Exit(0)
	},
}

