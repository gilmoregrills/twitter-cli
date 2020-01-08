package format

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/fatih/color"
)

// prints a tweet, calls itself when printing tweets which quote other tweets
// or are replies so that it can format threads/quote tweets properly
func PrintTweet(tweet *twitter.Tweet, indent int, client *twitter.Client) {
	if indent > 10 {
		return
	}
	indentStr := strings.Repeat(" ", indent)
	// fmt.Printf("%+v\n", tweet)
	color.Blue("%s@%s\n", indentStr, tweet.User.ScreenName)
	fmt.Printf("%s%+v\n", indentStr, tweet.Text)
	if tweet.QuotedStatus != nil {
		fmt.Printf("%squoted tweet:\n", indentStr)
		PrintTweet(tweet.QuotedStatus, indent+2, client)
	}
	if tweet.InReplyToStatusID != 0 {
		fmt.Printf("%sreplying to @%s\n", indentStr, tweet.InReplyToScreenName)
		previousTweet, _, err := client.Statuses.Show(tweet.InReplyToStatusID, nil)
		if err != nil {
			log.Fatal(err)
		}
		PrintTweet(previousTweet, indent+2, client)
	}
	color.Cyan("%slikes: %s, retweets: %s, replies: %s, quotes: %s\n", indentStr, strconv.Itoa(tweet.FavoriteCount), strconv.Itoa(tweet.RetweetCount), strconv.Itoa(tweet.ReplyCount), strconv.Itoa(tweet.QuoteCount))
}
