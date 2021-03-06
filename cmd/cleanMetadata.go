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
	"composer/service/clean"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "定时清理已经无用的Metadata",
	Run: func(cmd *cobra.Command, args []string) {
		cleanType := cmd.Flag("type").Value.String()
		runCleanMetadata(cleanType)
	},
}

func init() {
	cleanCmd.PersistentFlags().String("type", "metadata", "clean type dist or metadata")

	rootCmd.AddCommand(cleanCmd)
}

func runCleanMetadata(cleanType string) {
	if cleanType == "metadata" {
		clean.Metadata()
	}else if cleanType == "dist"{
		clean.Dist()
	}
}
