package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "youtube-cli",
	Short: "A command-line tool for interacting with the YouTube API",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("secrets", "", "path to your client secrets file")
	viper.BindPFlag("secrets", rootCmd.PersistentFlags().Lookup("secrets"))
	viper.BindEnv("secrets", "YOUTUBE_SECRETS")
}

func initConfig() {
	viper.AutomaticEnv()
}
