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
	// "time"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/fatih/color"
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

	// print the last 5 tweets
	// color.Yellow("fetching 5 most recent tweets to start with")
	// tweets, resp, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
	//    	Count: 5,
	// })
	// if err != nil {
	// 	log.Println(resp)
	// 	log.Fatal(err)
	// }
	// for i := 0; i < len(tweets); i++ {
	// 	printTweet(&tweets[i], 0, client)
	// }

	// get intial list and user.FriendsCount of those user is following
	// TODO: cache this list somewhere? don't want to generate every time
	//       we stream because we'll get rate limited, maybe I could check
	//       if the followercount has changed since last time and then if
	//       yes then update, that's a constant 1 api call every time stream runs not like 4+
	params := &twitter.FriendListParams{
		Cursor: 0,
		Count:  200,
	}
	friends, _, err := client.Friends.List(params)
	if err != nil {
		log.Fatal(err)
	}
	user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: viper.GetString("username"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// for user in list, get the IDStr and add it to a list
	following := make([]string, 0, user.FriendsCount)
	for friends.NextCursor != 0 {
		for user := range friends.Users {
			following = append(following, friends.Users[user].IDStr)
		}
		params.Cursor = friends.NextCursor
		friends, _, err = client.Friends.List(params)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(len(following))
	fmt.Println(following)

	// create a stream that follows that list of users
	// TODO: Filter this more so that it only returns tweets/RTs by that user
	streamParams := &twitter.StreamFilterParams{
		Follow:        following,
		FilterLevel:   viper.GetString("streaming_filter_level"),
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(streamParams)
	if err != nil {
		log.Fatal(err)
	}

	// let go-twitter handle the stream, executing the handler in demux.Tweet
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		//TODO: dont print if it's a tweet being RT'd by someone else etc
		if idStrInSlice(tweet.User.IDStr, following) {
			printTweet(tweet, 0, client)
			color.Yellow("---")
		}
	}
	go demux.HandleChan(stream.Messages)

	// handle closing the stream
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	color.Yellow("stopping stream...")
	stream.Stop()
}

func idStrInSlice(idStr string, users []string) bool {
	for _, user := range users {
		if user == idStr {
			return true
		}
	}
	return false
}

// prints a tweet, calls itself when printing tweets which quote other tweets
// or are replies so that it can format threads/quote tweets properly
// TODO: Make the formatting nicer
// TODO: Move to shared library ../lib/printers or something
func printTweet(tweet *twitter.Tweet, indent int, client *twitter.Client) {
	if indent > 10 {
		return
	}
	indentStr := strings.Repeat(" ", indent)
	// fmt.Printf("%+v\n", tweet)
	color.Blue("%s@%s\n", indentStr, tweet.User.ScreenName)
	fmt.Printf("%s%+v\n", indentStr, tweet.Text)
	if tweet.QuotedStatus != nil {
		fmt.Printf("%squoted tweet:\n", indentStr)
		printTweet(tweet.QuotedStatus, indent+2, client)
	}
	if tweet.InReplyToStatusID != 0 {
		fmt.Printf("%sreplying to @%s\n", indentStr, tweet.InReplyToScreenName)
		previousTweet, _, err := client.Statuses.Show(tweet.InReplyToStatusID, nil)
		if err != nil {
			log.Fatal(err)
		}
		printTweet(previousTweet, indent+2, client)
	}
	color.Cyan("%slikes: %s, retweets: %s, replies: %s, quotes: %s\n", indentStr, strconv.Itoa(tweet.FavoriteCount), strconv.Itoa(tweet.RetweetCount), strconv.Itoa(tweet.ReplyCount), strconv.Itoa(tweet.QuoteCount))
}

func init() {
	rootCmd.AddCommand(streamCmd)

	streamCmd.PersistentFlags().StringP("streaming_filter_level", "", "", "loaded from config")
	viper.BindPFlag("streaming_filter_level", streamCmd.PersistentFlags().Lookup("streaming_filter_level"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
