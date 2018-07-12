package cmd

import (
	"github.com/spf13/cobra"
)

const (
	deleteCommandUse = "delete"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(rmCmd)
}

var deleteCmd = &cobra.Command{
	Use:  deleteCommandUse,
	RunE: deleteExecute,
}

var rmCmd = &cobra.Command{
	Use:  "rm",
	RunE: deleteExecute,
}

func deleteExecute(_ *cobra.Command, _ []string) error {
	return client.Delete(path)
}
