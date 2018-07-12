package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getACLCmd)
}

var getACLCmd = &cobra.Command{
	Use:  "getacl",
	RunE: getACLExecute,
}

func getACLExecute(_ *cobra.Command, _ []string) error {
	value, err := client.GetACL(path)
	if err != nil {
		return err
	}

	out.PrintStringArray(value)

	return nil
}
