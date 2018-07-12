package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(existsCmd)
}

var existsCmd = &cobra.Command{
	Use:  "exists",
	RunE: existsExecute,
}

func existsExecute(_ *cobra.Command, _ []string) error {
	exists, err := client.Exists(path)
	if err != nil {
		return err
	}

	fmt.Println(exists)

	return nil
}
