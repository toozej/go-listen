package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/toozej/go-listen/pkg/config"
	"github.com/toozej/go-listen/pkg/man"
	"github.com/toozej/go-listen/pkg/version"
)

var conf config.Config

var rootCmd = &cobra.Command{
	Use:              "go-listen",
	Short:            "Spotify playlist management tool",
	Long:             `go-listen is a web application that allows users to search for artists and automatically add their top 5 songs to designated "incoming" playlists on Spotify.`,
	Args:             cobra.ExactArgs(0),
	PersistentPreRun: rootCmdPreRun,
	Run:              rootCmdRun,
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	// Show help when no subcommand is provided
	if err := cmd.Help(); err != nil {
		log.WithError(err).Error("Failed to show help")
	}
}

func rootCmdPreRun(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// Load configuration with debug flag
	conf = config.GetEnvVars(debug)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	// create rootCmd-level flags
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug-level logging")

	// add sub-commands
	rootCmd.AddCommand(
		man.NewManCmd(),
		version.Command(),
	)
}
