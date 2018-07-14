package cmd

import (
	"github.com/spf13/cobra"
)

const (
	getCommandUse = "get"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   getCommandUse,
	Short: "Get the value of the specified znode",
	RunE:  getExecute,
}

func getExecute(_ *cobra.Command, _ []string) error {
	value, err := client.Get(path)
	if err != nil {
		return err
	}

	out.Printf("%s", value)

	return nil
}
