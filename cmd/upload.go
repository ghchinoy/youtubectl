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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a video to YouTube",
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		category, _ := cmd.Flags().GetString("category")
		keywords, _ := cmd.Flags().GetString("keywords")
		privacy, _ := cmd.Flags().GetString("privacy")
		uploadVideo(filename, title, description, category, keywords, privacy)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().String("filename", "", "The video file to upload")
	uploadCmd.Flags().String("title", "Test Title", "The video title")
	uploadCmd.Flags().String("description", "Test Description", "The video description")
	uploadCmd.Flags().String("category", "22", "The video category")
	uploadCmd.Flags().String("keywords", "", "A comma-separated list of video keywords")
	uploadCmd.Flags().String("privacy", "unlisted", "The video privacy status")
	uploadCmd.MarkFlagRequired("filename")
}

func uploadVideo(filename, title, description, category, keywords, privacy string) {
	ctx := context.Background()

	secretsFile := viper.GetString("secrets")
	if secretsFile == "" {
		log.Fatalf("Please provide the path to your client secrets file with the --secrets flag")
	}

	b, err := ioutil.ReadFile(secretsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(ctx, config)
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: privacy},
	}

	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}
	defer file.Close()

	response, err := call.Media(file).Do()
	if err != nil {
		log.Fatalf("Error making API call: %v", err)
	}

	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
}
