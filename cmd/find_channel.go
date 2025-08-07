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

var findChannelCmd = &cobra.Command{
	Use:   "find-channel",
	Short: "Find a channel ID by username",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		findChannel(username)
	},
}

func init() {
	rootCmd.AddCommand(findChannelCmd)
	findChannelCmd.Flags().String("username", "", "The username to search for")
	findChannelCmd.MarkFlagRequired("username")
}

func findChannel(username string) {
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

	call := service.Search.List([]string{"snippet"}).Q(username).Type("channel")
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call to search for channel: %v", err)
	}

	if len(response.Items) == 0 {
		fmt.Println("Could not find a channel with that username.")
		return
	}

	channel := response.Items[0]

	fmt.Printf("Found channel:\n")
	fmt.Printf("  Title: %s\n", channel.Snippet.Title)
	fmt.Printf("  Channel ID: %s\n", channel.Snippet.ChannelId)
	fmt.Printf("  Description: %s\n", channel.Snippet.Description)
}
