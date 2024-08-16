package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "integration-git",
	Short: "Lightweight UI, that presents the findings from the latest scan and prompt the user to review pending identifications.",
}

func Execute() error {
	return rootCmd.Execute()
}
