package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/laetho/metagraf/internal/pkg/params/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
)

func init() {
	createCmd.AddCommand(createRefCmd)
	createRefCmd.Flags().StringVarP(&Namespace, "namespace", "n","","namespace to fetch template form")
	createRefCmd.Flags().StringVarP(&Template, "template", "t", "metagraf-refdoc.md", "name of ConfigMap for go template")
	createRefCmd.Flags().StringVarP(&Suffix, "suffix", "s", ".html", "file suffix of the generated content")
	createRefCmd.Flags().StringVarP(&params.RefTemplateFile, "templatefile", "f", "", "file path to go template to use when creating reference doucment")
	createRefCmd.Flags().StringVarP(&params.RefTemplateOutputFile, "outputfile", "O", "", "output file, may include leading path")
}

var createRefCmd = &cobra.Command{
	Use:   "ref <metaGraf>",
	Short: "create ref document from metaGraf specification",
	Long:  MGBanner + `create ref`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(StrMissingCollection)
			return
		}

		mg := metagraf.Parse(args[0])
		FlagPassingHack()
		modules.GenRef(&mg)
	},
}