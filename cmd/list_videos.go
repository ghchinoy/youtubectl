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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var listVideosCmd = &cobra.Command{
	Use:   "list-videos",
	Short: "List videos for a specific channel",
	Run: func(cmd *cobra.Command, args []string) {
		channelID, _ := cmd.Flags().GetString("channelid")
		detailed, _ := cmd.Flags().GetBool("detailed")
		limit, _ := cmd.Flags().GetInt64("limit")
		listVideos(channelID, detailed, limit)
	},
}

func init() {
	rootCmd.AddCommand(listVideosCmd)
	listVideosCmd.Flags().String("channelid", "", "The YouTube channel ID to list videos for")
	listVideosCmd.Flags().Bool("detailed", false, "Display detailed video information (statistics, status, etc.)")
	listVideosCmd.Flags().Int64("limit", 10, "The maximum number of results to display per page")
	listVideosCmd.MarkFlagRequired("channelid")
}

func listVideos(channelID string, detailed bool, limit int64) {
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

	channelResponse, err := service.Channels.List([]string{"contentDetails"}).Id(channelID).Do()
	if err != nil || len(channelResponse.Items) == 0 {
		log.Fatalf("Error finding channel or channel not found: %v", err)
	}
	uploadsPlaylistID := channelResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads

	nextPageToken := ""
	for {
		playlistCall := service.PlaylistItems.List([]string{"snippet"}).
			PlaylistId(uploadsPlaylistID).
			MaxResults(limit).
			PageToken(nextPageToken)

		playlistItemsResponse, err := playlistCall.Do()
		if err != nil {
			log.Fatalf("Error getting playlist items: %v", err)
		}

		if len(playlistItemsResponse.Items) == 0 {
			fmt.Println("No videos found on this channel.")
			return
		}

		var videoIDs []string
		for _, item := range playlistItemsResponse.Items {
			videoIDs = append(videoIDs, item.Snippet.ResourceId.VideoId)
		}

	
videosCall := service.Videos.List([]string{"snippet", "statistics", "status"}).Id(strings.Join(videoIDs, ","))
	
videosResponse, err := videosCall.Do()
		if err != nil {
			log.Fatalf("Error getting video details: %v", err)
		}

		fmt.Printf("Displaying videos for channel %s:\n\n", channelID)
		for _, video := range videosResponse.Items {
			fmt.Printf("Title: %s\n", video.Snippet.Title)
			fmt.Printf("Video ID: %s\n", video.Id)
			fmt.Printf("Published At: %s\n", video.Snippet.PublishedAt)
			if detailed {
				fmt.Printf("View Count: %d\n", video.Statistics.ViewCount)
				fmt.Printf("Like Count: %d\n", video.Statistics.LikeCount)
				fmt.Printf("Comment Count: %d\n", video.Statistics.CommentCount)
				fmt.Printf("Privacy: %s\n", video.Status.PrivacyStatus)
			}
			fmt.Println("--------------------------------------------------")
		}

		nextPageToken = playlistItemsResponse.NextPageToken
		if nextPageToken == "" {
			break
		}

		fmt.Print("Show more? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			break
		}
	}
}