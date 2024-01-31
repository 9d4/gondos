package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gondos",
	Short: "Gondos! Hello World",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	// Register commands here
	rootCmd.AddCommand(newServeCmd().cmd)
}
