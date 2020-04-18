package cmd

import (
	_import "composer/service/import"
	"github.com/spf13/cobra"
)

var importDist = &cobra.Command{
	Use:   "import:dist",
	Short: "将dist云存储中的文件数据导入redis",
	Run: func(cmd *cobra.Command, args []string) {
		runImportDist(cmd.Flag("nextPage").Value.String())
	},
}

func init() {
	importDist.PersistentFlags().String("nextPage", "", "nextPage")
	rootCmd.AddCommand(importDist)
}

func runImportDist(nextPage string) {

	_import.ImportDist(nextPage)
}
