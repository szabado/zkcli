package cmd

import (
	"github.com/spf13/cobra"
)

const (
	existsCommandUse = "exists"
)

func init() {
	rootCmd.AddCommand(existsCmd)
}

var existsCmd = &cobra.Command{
	Use:   existsCommandUse,
	Short: "Check if the specified znode exists",
	RunE:  existsExecute,
}

func existsExecute(_ *cobra.Command, _ []string) error {
	exists, err := client.Exists(path)
	if err != nil {
		return err
	}

	out.Printf("%v", exists)

	return nil
}
