package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getACLCmd)
}

var getACLCmd = &cobra.Command{
	Use:  "getacl",
	Short: "Get the ACL associated with a znode",
	RunE: getACLExecute,
}

func getACLExecute(_ *cobra.Command, _ []string) error {
	value, err := client.GetACL(path)
	if err != nil {
		return err
	}

	out.PrintArray(value)

	return nil
}
