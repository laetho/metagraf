/*
Copyright 2018 The metaGraf Authors

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
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	log "k8s.io/klog"
)

func init() {
	RootCmd.AddCommand(devCmd)
	devCmd.PersistentFlags().BoolVar(&Output, "output", false, "also output objects")
	devCmd.PersistentFlags().BoolVar(&Dryrun, "dryrun", false, "do not create objects, only output")
	devCmd.PersistentFlags().StringVarP(&Format, "format", "o", "json", "specify json or yaml, json id default")

	devCmd.AddCommand(devCmdUp)
	devCmdUp.Flags().StringVarP(&params.NameSpace, "namespace", "n", "", "namespace to work on, if not supplied it will use current active namespace.")
	devCmdUp.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	devCmdUp.Flags().StringVar(&params.PropertiesFile, "cvfile", "", "Property file with component configuration values. Can be generated with \"mg generate properties\" command.)")
	devCmdUp.Flags().StringVar(&OName, "name", "", "Overrides name of application.")
	devCmdUp.Flags().StringVarP(&Registry, "registry", "r", viper.GetString("registry"), "Specify container registry host")
	devCmdUp.Flags().StringVarP(&params.OutputImagestream, "istream", "i", "", "specify if you want to output to another imagestream than the component name")
	devCmdUp.Flags().StringVarP(&Context, "context", "c", "/", "Application contextroot. (\"/<context>\"). Used when creating Route object.")
	devCmdUp.Flags().BoolVarP(&CreateGlobals, "globals", "g", false, "Override default behavior and force creation of global secrets. Will not overwrite existing ones.")
	devCmdUp.Flags().BoolVar(&params.CreateSecrets, "create-secrets", false, "Creates empty secrets referenced in metagraf specification. Needs to be manually filled with values.")
	devCmdUp.Flags().BoolVar(&params.ServiceMonitor, "service-monitor", false, "Set flag to also create a ServiceMonitor resource. Requires a cluster with the prometheus-operator.")
	devCmdUp.Flags().Int32Var(&params.ServiceMonitorPort, "service-monitor-port", params.ServiceMonitorPort, "Set Service port to scrape in ServiceMonitor.")
	devCmdUp.Flags().StringVar(&params.ServiceMonitorOperatorName, "service-monitor-operator-name", params.ServiceMonitorOperatorName, "Name of prometheus-operator instance to create ServiceMonitor for.")


	devCmd.AddCommand(devCmdDown)
	devCmdDown.Flags().StringVarP(&params.NameSpace, "namespace", "n", "", "namespace to work on, if not supplied it will use current active namespace.")
	devCmdDown.Flags().BoolVar(&params.Everything, "everything", false, "Delete all resources and artifacts generated from mg dev up.")
	devCmdDown.Flags().StringVar(&OName, "name", "", "Overrides name of application.")

	devCmd.AddCommand(devCmdBuild)
	devCmdBuild.Flags().StringVar(&params.SourceRef, "ref", "master", "Specify the git ref or branch ref to build.")
	devCmdBuild.Flags().StringVarP(&params.NameSpace, "namespace", "n", "", "namespace to work on, if not supplied it will use current active namespace.")
	devCmdBuild.Flags().BoolVar(&params.LocalBuild, "local", false, "Builds application from src in current (.) direcotry.")
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "dev subcommands",
	Long:  `dev subcommands`,
}

var devCmdUp = &cobra.Command{
	Use:   "up <metagraf.json>",
	Short: "creates the required component resources.",
	Long:  `sets up the required resources to test the component in the platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		requireMetagraf(args)
		requireNamespace()

		FlagPassingHack()
		modules.NameSpace = params.NameSpace

		devUp(args[0])
	},
}

var devCmdDown = &cobra.Command{
	Use:   "down <metagraf.json>",
	Short: "deletes component resources",
	Long:  `dev subcommands`,
	Run: func(cmd *cobra.Command, args []string) {
		requireMetagraf(args)
		requireNamespace()

		FlagPassingHack()
		modules.NameSpace = params.NameSpace

		devDown(args[0])
	},
}

var devCmdBuild = &cobra.Command{
	Use:   "build <metagraf.json>",
	Short: "build the container from generated BuildConfig",
	Long:  `dev subcommands`,
	Run: func(cmd *cobra.Command, args []string) {
		requireMetagraf(args)
		requireNamespace()

		mg := metagraf.Parse(args[0])
		bc := mg.Name(OName, Version)

		FlagPassingHack()
		modules.NameSpace = params.NameSpace

		// Remove RepSecRef from generated BuildConfig if --local argument is provided.
		if params.LocalBuild && len(mg.Spec.RepSecRef) > 0 {
			mg.Spec.RepSecRef = ""
		}

		modules.GenImageStream(&mg, params.NameSpace)
		modules.GenBuildConfig(&mg)

		path, err := exec.LookPath("oc")
		if err != nil {
			log.Fatal(err)
		}
		arg := []string{"start-build", bc,"-n", params.NameSpace, "--follow"}
		if params.LocalBuild {
			arg = append(arg, "--from-file=.")
		}
		c := exec.Command(path, arg...)
		stdout, err := c.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		c.Env = append(os.Environ())
		c.Start()

		buf := bufio.NewReader(stdout)
		line := []byte{}
		for err == nil {
			line, _,err = buf.ReadLine()
			fmt.Println(string(line))
		}
	},
}

func devUp(mgf string) {
	mg := metagraf.Parse(mgf)
	modules.Variables = GetCmdProperties(mg.GetProperties())
	log.V(2).Info("Current MGProperties: ", modules.Variables)

	modules.PullPolicy = corev1.PullAlways
	modules.GenSecrets(&mg)
	modules.GenConfigMaps(&mg)
	modules.GenDeploymentConfig(&mg)
	modules.GenService(&mg)
	modules.GenRoute(&mg)

}

func devDown(mgf string) {
	mg := metagraf.Parse(mgf)
	basename := modules.Name(&mg)

	modules.DeleteRoute(basename)
	modules.DeleteService(basename)
	modules.DeleteServiceMonitor(basename)
	modules.DeleteDeploymentConfig(basename)
	modules.DeleteBuildConfig(basename)
	modules.DeleteConfigMaps(&mg)
	modules.DeleteImageStream(basename)

	if params.Everything {
		modules.DeleteSecrets(&mg)
	}
}

func checkForOC() bool {
	_, err := exec.LookPath("oc")
	if err != nil {
		return false
	}
	return true
}
