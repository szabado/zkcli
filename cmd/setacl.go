package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

const (
	setAclCommandUse = "setacl"
)

func init() {
	rootCmd.AddCommand(setACLCmd)
}

var setACLCmd = &cobra.Command{
	Use:   setAclCommandUse + " <path> <acl>",
	Short: "Set the ACL of the specified znode",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 2 {
			aclstr = args[1]
		} else {
			var err error
			data, err := io.ReadAll(stdin)
			if err != nil {
				return err
			}
			aclstr = string(data)
		}
		return nil
	},
	RunE: setACLExecute,
}

func setACLExecute(_ *cobra.Command, _ []string) error {
	result, err := client.SetACL(path, aclstr, force)
	if err != nil {
		return err
	}

	out.Printf("Set %+v", result)

	return nil
}
