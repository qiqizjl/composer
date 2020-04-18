package cmd

import (
	"composer/service/clean"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCloudMetaCmd = &cobra.Command{
	Use:   "clean:cloud:metadata",
	Short: "清理远程存在本地却不存在的文件",
	Run: func(cmd *cobra.Command, args []string) {
		runCleanCloudMeta()
	},
}

func init() {
	rootCmd.AddCommand(cleanCloudMetaCmd)
}

func runCleanCloudMeta() {
	clean.CloudMeta()
}
