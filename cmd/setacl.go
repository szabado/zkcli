package cmd

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setACLCmd)
}

var setACLCmd = &cobra.Command{
	Use: "setacl",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			aclstr = args[1]
		} else {
			var err error
			data, err := ioutil.ReadAll(os.Stdin)
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

	log.Infof("Set %+v", result)

	return nil
}
