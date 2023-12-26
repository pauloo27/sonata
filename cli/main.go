package main

import (
	"os"

	"github.com/pauloo27/sonata/cli/command"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sonata",
	Short: "Sleepy HTTP Client",
	Long:  "Run HTTP requests and etc",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(command.Version)
	rootCmd.AddCommand(command.New)
	rootCmd.AddCommand(command.Run)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
