package cmd

import (
	"github.com/spf13/cobra"
)

const (
	deleteCommandUse = "delete"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:  deleteCommandUse,
	Short: "Delete the specified znode",
	Aliases:[]string{"rm"},
	RunE: deleteExecute,
}

func deleteExecute(_ *cobra.Command, _ []string) error {
	return client.Delete(path)
}
