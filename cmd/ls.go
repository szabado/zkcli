package cmd

import (
	"github.com/spf13/cobra"
)

const (
	lsCommandUse = "ls"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   lsCommandUse,
	Short: "Get the children of the specified znode",
	RunE:  lsExecute,
}

func lsExecute(_ *cobra.Command, _ []string) error {
	children, err := client.Children(path)
	if err != nil {
		return err
	}

	out.PrintArray(children)

	return nil
}
