package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var channelInfoCmd = &cobra.Command{
	Use:   "channel-info",
	Short: "Get channel information by video ID",
	Run: func(cmd *cobra.Command, args []string) {
		videoID, _ := cmd.Flags().GetString("videoid")
		getChannelInfoByVideoID(videoID)
	},
}

func init() {
	rootCmd.AddCommand(channelInfoCmd)
	channelInfoCmd.Flags().String("videoid", "", "The ID of a video on the channel")
	channelInfoCmd.MarkFlagRequired("videoid")
}

func getChannelInfoByVideoID(videoID string) {
	ctx := context.Background()

	secretsFile := viper.GetString("secrets")
	if secretsFile == "" {
		log.Fatalf("Please provide the path to your client secrets file with the --secrets flag or YOUTUBE_SECRETS env var")
	}

	b, err := ioutil.ReadFile(secretsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(ctx, config)
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	videoResponse, err := service.Videos.List([]string{"snippet"}).Id(videoID).Do()
	if err != nil {
		log.Fatalf("Error getting video info: %v", err)
	}

	if len(videoResponse.Items) == 0 {
		fmt.Println("Could not find a video with that ID.")
		return
	}

	channelID := videoResponse.Items[0].Snippet.ChannelId

	channelResponse, err := service.Channels.List([]string{"snippet", "contentDetails", "statistics"}).Id(channelID).Do()
	if err != nil {
		log.Fatalf("Error getting channel info: %v", err)
	}

	if len(channelResponse.Items) == 0 {
		fmt.Println("Could not find a channel associated with that video.")
		return
	}

	channel := channelResponse.Items[0]

	fmt.Printf("Channel ID: %s\n", channel.Id)
	fmt.Printf("Title: %s\n", channel.Snippet.Title)
	fmt.Printf("Description: %s\n", channel.Snippet.Description)
	fmt.Printf("View Count: %d\n", channel.Statistics.ViewCount)
	fmt.Printf("Subscriber Count: %d\n", channel.Statistics.SubscriberCount)
	fmt.Printf("Video Count: %d\n", channel.Statistics.VideoCount)
}
