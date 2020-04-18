/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	_ "composer/file"
	"composer/service/queue"
	"github.com/spf13/cobra"
	"net/http"
	_ "net/http/pprof"
)

// spiderCmd represents the spider command
var spiderCmd = &cobra.Command{
	Use:   "spider",
	Short: "抓取Composer数据",
	Run: func(cmd *cobra.Command, args []string) {
		runSpider()
	},
}

func init() {

	rootCmd.AddCommand(spiderCmd)
}

func runSpider() {
	go queue.Package(1)
	for i := 0; i < 12; i++ {
		go queue.Provider(i)
	}
	for i := 0; i < 30; i++ {
		go queue.PackageHash(i)
	}
	for i := 0; i < 50; i++ {
		go queue.Dist(i)
	}
	// 启动debug服务
	http.ListenAndServe("0.0.0.0:8080", nil)

}
