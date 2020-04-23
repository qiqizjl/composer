package cmd

import (
	"composer/file"
	"composer/service/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var sumDistSize = &cobra.Command{
	Use:   "sum:dist",
	Short: "统计每个包使用空间大小",
	Run: func(cmd *cobra.Command, args []string) {
		runSumDistSize(cmd.Flag("nextPage").Value.String())
	},
}

func init() {
	sumDistSize.PersistentFlags().String("nextPage", "", "nextPage")
	rootCmd.AddCommand(sumDistSize)
}

func runSumDistSize(nextPage string) {
	if nextPage == "" {
		redis.ClearDistSize()
	}
	result, err := file.DistFile.ListFile(nextPage, 100000)
	if err != nil {
		logrus.Errorln("request error", nextPage, err.Error())
		return
	}
	for item := range result {
		logrus.Infoln("update ", item.Remark)
		s := strings.Split(item.Key, "/")
		packageName := s[0] + "/" + s[1]
		redis.AddDistSize(packageName,item.Size)
	}
}
