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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/fatih/color"
	"github.com/gilmoregrills/twitter-cli/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream your twitter feed.",
	Long: `Fetches a few recent tweets, then streams incoming tweets on top.
TODO: Detail how to interact with the feed here.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow("streaming %s's timeline", viper.GetString("username"))
		stream()
	},
}

func stream() {
	config := oauth1.NewConfig(viper.GetString("consumer_token"), viper.GetString("consumer_token_secret"))
	token := oauth1.NewToken(viper.GetString("access_token"), viper.GetString("access_token_secret"))
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	color.Yellow("fetching 5 most recent tweets to start with")
	tweets, resp, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 5,
	})
	if err != nil {
		log.Println(resp)
		log.Fatal(err)
	}
	for i := 0; i < len(tweets); i++ {
		format.PrintTweet(&tweets[i], 0, client)
	}
	// followers, resp, err := client.Followers.List(&twitter.FollowerListParams{})
	// for follower in followers, get the IDStr and add it to a list pls
	// fmt.Printf("%+v\n", followers)

	params := &twitter.StreamUserParams{
		With:          "followings",
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.User(params)
	if err != nil {
		log.Fatal(err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		format.PrintTweet(tweet, 0, client)
		color.Yellow("---")
	}

	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	color.Yellow("stopping stream...")
	stream.Stop()
}

func init() {
	rootCmd.AddCommand(streamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
