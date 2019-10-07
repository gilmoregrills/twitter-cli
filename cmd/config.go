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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Prints config file contents, and allows updating.",
	Long: `TODO: Detail available options here.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
		fmt.Println("username is: ", viper.GetString("username"))
		fmt.Println("access token is: ", viper.GetString("access_token"))
		fmt.Println("access token secret is: ", viper.GetString("access_token_secret"))
		fmt.Println("consumer token is: ", viper.GetString("consumer_token"))
		fmt.Println("consumer token secret is: ", viper.GetString("consumer_token_secret"))
		fmt.Println("streaming filter level is: ", viper.GetString("streaming_filter_level"))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
