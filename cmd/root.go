/*
Copyright © 2019 gilmoregrills

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
    "os"
    "github.com/spf13/cobra"

    homedir "github.com/mitchellh/go-homedir"
    "github.com/spf13/viper"

)


var cfgFile string


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "twitter-cli",
    Short: "Streams your twitter feed in the command line, more coming",
    Long: `It seemed like a fun exercise at the time. Some problems to
solve include:
  - How to format tweets in a nice presentable way
  - How to handle images!!
  - Should there be a way to select tweets and view replies?
  - How should we handle responses and stuff?`,
    // Run: func(cmd *cobra.Command, args []string) {
    //     api := anaconda.NewTwitterApiWithCredentials(viper.GetString("access_token"), viper.GetString("access_token_secret"), viper.GetString("consumer_token"), viper.GetString("consumer_token_secret"))
    // }
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.twitter-cli/config)")
    // Your twitter username
    rootCmd.PersistentFlags().StringP("username", "", "", "loaded from config")
    viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
    // Your twitter credentials
    rootCmd.PersistentFlags().StringP("access_token", "", "", "loaded from config")
    viper.BindPFlag("access_token", rootCmd.PersistentFlags().Lookup("access_token"))
    rootCmd.PersistentFlags().StringP("access_token_secret", "", "", "loaded from config")
    viper.BindPFlag("access_token_secret", rootCmd.PersistentFlags().Lookup("access_token_secret"))
    rootCmd.PersistentFlags().StringP("consumer_token", "", "", "loaded from config")
    viper.BindPFlag("consumer_token", rootCmd.PersistentFlags().Lookup("consumer_token"))
    rootCmd.PersistentFlags().StringP("consumer_token_secret", "", "", "loaded from config")
    viper.BindPFlag("consumer_token_secret", rootCmd.PersistentFlags().Lookup("consumer_token_secret"))

     // Cobra also supports local flags, which will only run
      // when this action is called directly.
    rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Find home directory.
        home, err := homedir.Dir()
        if err != nil {
        fmt.Println(err)
        os.Exit(1)
        }

        // Search config in home directory with name ".twitter-cli" (without extension).
        viper.AddConfigPath(home)
        viper.SetConfigName(".twitter-cli")
    }

    viper.AutomaticEnv() // read in environment variables that match

    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    }
}

