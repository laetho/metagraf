package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	params "metagraf/internal/pkg/params/params"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"os"
)

func init() {
	RootCmd.AddCommand(InspectCmd)
	InspectCmd.Flags().BoolVar(&Enforce, "enforce", false, "Enforce findings, defaults to false and informs only.")
	InspectCmd.AddCommand(InspectPropertiesCmd)
	InspectPropertiesCmd.Flags().StringVar(&params.PropertiesFile, "cvfile", "", "File with component configuration values. (key=value pairs)")
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
		modules.Variables = GetCmdProperties(mg.GetProperties())
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
			log.Infof("Active project is: %v", viper.Get("namespace"))
			log.Errorf("Missing path to metaGraf specification")
			os.Exit(1)
		}

		if len(args) < 2 {
			log.Errorf("Missing path to properties file")
			os.Exit(1)
		}

		mg := metagraf.Parse(args[0])
		params.PropertiesFile = args[1]
		modules.Variables = GetCmdProperties(mg.GetProperties())

		if !ValidateProperties(modules.Variables) {
			os.Exit(1)
		} else {
			fmt.Printf("The %v configuration is valid for this metaGraf specification.\n", params.PropertiesFile)
		}
	},
}

// Check if all required MGProperty structs in MGProperties has a value.
// Returns true if they are valid, false if they are invalid.
func ValidateProperties(mgprops metagraf.MGProperties) bool {
	reqvars := mgprops.GetRequired().SourceKeyMap(true)

	fail := false
	for key, _ := range reqvars {
		property := mgprops[key]
		if len(property.Value) == 0 {
			fail = true
			fmt.Printf("Required parameter %v does not have a value.\n", property.MGKey())
		}
	}

	if fail {
		return false
	}
	return true
}
