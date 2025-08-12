// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List channel information by ID or username",
	Run: func(cmd *cobra.Command, args []string) {
		query, _ := cmd.Flags().GetString("query")
		listChannels(query)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("query", "", "The channel ID or username to get information for")
	listCmd.MarkFlagRequired("query")
}

func listChannels(query string) {
	ctx := context.Background()

	secretsFile := viper.GetString("secrets")
	if secretsFile == "" {
		log.Fatalf("Please provide the path to your client secrets file with the --secrets flag")
	}

	b, err := ioutil.ReadFile(secretsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope, youtube.YoutubeUploadScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(ctx, config)
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// First, try to get the channel by ID
	call := service.Channels.List([]string{"snippet", "contentDetails", "statistics"}).Id(query)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call: %v", err)
	}

	// If that fails, try searching for the channel by username
	if len(response.Items) == 0 {
		searchCall := service.Search.List([]string{"snippet"}).Q(query).Type("channel")
		searchResponse, err := searchCall.Do()
		if err != nil {
			log.Fatalf("Error making API call to search for channel: %v", err)
		}

		if len(searchResponse.Items) == 0 {
			fmt.Println("Could not find a channel with that ID or username.")
			return
		}

		// Now, get the full channel details for the first search result
		channelID := searchResponse.Items[0].Snippet.ChannelId
		call = service.Channels.List([]string{"snippet", "contentDetails", "statistics"}).Id(channelID)
		response, err = call.Do()
		if err != nil {
			log.Fatalf("Error making API call: %v", err)
		}
	}

	channel := response.Items[0]

	fmt.Printf("Channel ID: %s\n", channel.Id)
	fmt.Printf("Title: %s\n", channel.Snippet.Title)
	fmt.Printf("Description: %s\n", channel.Snippet.Description)
	fmt.Printf("View Count: %d\n", channel.Statistics.ViewCount)
	fmt.Printf("Subscriber Count: %d\n", channel.Statistics.SubscriberCount)
	fmt.Printf("Video Count: %d\n", channel.Statistics.VideoCount)
}
