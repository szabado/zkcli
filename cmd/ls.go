package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:  "ls",
	RunE: lsExecute,
}

func lsExecute(_ *cobra.Command, _ []string) error {
	children, err := client.Children(path)
	if err != nil {
		return err
	}

	out.PrintStringArray(children)

	return nil
}
