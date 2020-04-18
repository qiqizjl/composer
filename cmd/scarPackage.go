package cmd

import (
	"composer/service/clean"
	"composer/service/redis"
	"github.com/spf13/cobra"
	"time"
)

// scarPackageCmd represents the clean command
var scarPackageCmd = &cobra.Command{
	Use:   "scarPackage",
	Short: "扫描package",
	Run: func(cmd *cobra.Command, args []string) {
		runScarPackage()
	},
}

func init() {
	rootCmd.AddCommand(scarPackageCmd)
}

func runScarPackage() {
	startTime := time.Now().Unix()
	err := clean.UpdateMetadataTime(true)
	if err == nil {
		redis.SetUpdateTime(startTime)
	}
}
