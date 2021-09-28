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
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	log "k8s.io/klog"
)

var wg sync.WaitGroup

// Type for events emitted by filteredFileWatcher.
type WatchEvent int

// Enum values for WatchEvent
const (
	BuildAndDeploy WatchEvent = iota
	Deploy
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
	devCmdBuild.Flags().StringVar(&params.SourceRef, "ref", "", "Specify the git ref or branch ref to build.")
	devCmdBuild.Flags().StringVarP(&params.NameSpace, "namespace", "n", "", "namespace to work on, if not supplied it will use current active namespace.")
	devCmdBuild.Flags().BoolVar(&params.LocalBuild, "local", false, "Builds application from src in current (.) direcotry.")
	devCmdBuild.Flags().StringSliceVar(&params.BuildParams, "buildparams", []string{}, "Slice of key=value pairs, seperated by , to override or append build params used by underlying build mechanism")

	devCmd.AddCommand(devCmdWatch)
	devCmdWatch.Flags().StringVarP(&params.NameSpace, "namespace", "n", "", "namespace to work on, if not supplied it will use current active namespace.")
	devCmdWatch.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	devCmdWatch.Flags().StringSliceVar(&IgnoredPaths, "ignore-paths", []string{}, "List of paths to ignore when watching file changes, seperated by \",\".")
	devCmdWatch.Flags().StringVar(&params.PropertiesFile, "cvfile", "", "Property file with component configuration values. Can be generated with \"mg generate properties\" command.)")
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "dev subcommands",
	Long:  `dev subcommands`,
}

var devCmdWatch = &cobra.Command{
	Use:   "watch <metagraf.json>",
	Short: "watches for local filechanges and rebuilds or redeploys the component.",
	Long:  `Inspired by skaffold.dev. Watches for local file changes and rebuilds and redeploys.`,
	Run: func(cmd *cobra.Command, args []string) {
		requireMetagraf(args)
		requireNamespace()

		mg := metagraf.Parse(args[0])
		bc := mg.Name(OName, Version)

		FlagPassingHack()
		modules.NameSpace = params.NameSpace

		// Crate a buffered channel with room for one event.
		chEvents := make(chan WatchEvent, 1)
		// Channel to indicate if we are in processing state.
		chProcessing := make(chan bool, 1)

		wg.Add(2)
		defer wg.Wait()
		go filteredFileWatcher(chEvents, chProcessing)

		for {
			select {
			case command := <-chEvents:
				chProcessing<-true
				switch command {
				case BuildAndDeploy:
					buildGenerate(&mg, params.NameSpace, true)
					err := s2ibuild(bc, params.NameSpace, true)
					if err != nil {
						log.Fatalf("Unable to build: %v ", err)
					}
					devUp(args[0])
				case Deploy:
					devUp(args[0])
				}
				chProcessing<-false
			}
		}
	},
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

		buildGenerate(&mg, params.NameSpace, params.LocalBuild)
		err := s2ibuild(bc, params.NameSpace, params.LocalBuild)
		if err != nil {
			log.Fatalf("Build failed: %v", err)
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

func buildGenerate(mg *metagraf.MetaGraf, ns string, local bool) {
	// Remove RepSecRef from generated BuildConfig if --local argument is provided.
	if local && len(mg.Spec.RepSecRef) > 0 {
		mg.Spec.RepSecRef = ""
	}

	modules.GenImageStream(mg, ns)
	modules.GenBuildConfig(mg)
}

func s2ibuild(bc string, ns string, local bool) error {
	path, err := exec.LookPath("oc")
	if err != nil {
		return err
	}

	arg := []string{"start-build", bc, "-n", ns, "--follow"}

	if local {
		arg = append(arg, "--from-dir=.")
		arg = append(arg, "--exclude=''")
	}

	c := exec.Command(path, arg...)
	stdout, pipeerr := c.StdoutPipe()
	if pipeerr != nil {
		return pipeerr
	}
	c.Env = append(os.Environ())
	cmderr := c.Start()
	if cmderr != nil {
		return cmderr
	}

	buf := bufio.NewReader(stdout)
	line := []byte{}
	for err == nil {
		line, _, err = buf.ReadLine()
		fmt.Println(string(line))
	}
	return nil
}

// Watches for events in ./ using fsnotify and writes WatchEvent's to typed channel.
// Reads processing state from a bool channel. If true, skip writing commands.
func filteredFileWatcher(chEvents chan<- WatchEvent, chProcessing <-chan bool) {

	done := make(chan bool)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	watchDir := func(path string, info os.FileInfo, err error) error {
		for _,ignored := range IgnoredPaths {
			if strings.Contains(path, ignored) {
				return nil
			}
		}
		if info.Mode().IsDir() {
			return watcher.Add(path)
		}
		return nil
	}

	if err := filepath.Walk("./", watchDir); err != nil {
		log.Fatal(err)
	}

	go func() {
		var skipevents bool = false
		for {
			select {
			case event := <-watcher.Events:
				// Check to see if we got a new chProcessing state to handle
				if len(chProcessing) >= 1 {
					skipevents = <-chProcessing
				}
				if skipevents {
					// Pop of the Event and continue.
					<-watcher.Events
					continue
				}

				// Discard temporary files from editors
				if strings.Contains(event.Name, "~") {
					continue
				}
				if len(event.Name) == 0 {
					continue
				}

				switch event.Op.String() {
				case "WRITE":
					log.V(2).Infof("Got event on: %v, Type: %v", event.Name, event.Op)
					// Redeploy if configuration input changed.
					if strings.Contains(event.Name, params.PropertiesFile) {
						log.Info("Properties file changed, redeploying.")
						chEvents <- Deploy
					} else if strings.Contains(event.Name, "metagraf.json") {
						log.Info("metaGraf Specification Changed")
						chEvents <- BuildAndDeploy
					}else {
						// Build and Redeploy
						chEvents <- BuildAndDeploy
					}
				}
			case err := <-watcher.Errors:
				fmt.Println("ERROR:", err)
			default:
			}
		}
	}()
	<-done
}
