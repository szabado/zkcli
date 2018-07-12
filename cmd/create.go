package cmd

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	createCommandUse = "create"
)

var (
	aclstr string
	acls   string
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVar(&acls, "acls", "31", "optional, csv list [1|,2|,4|,8|,16|,31]")
}

var createCmd = &cobra.Command{
	Use: createCommandUse,
	Short: "Create the specified znode",
	// Don't have preruns here, creater calls createExecute and does no setup
	RunE: createExecute,
}

func createExecute(_ *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.Errorf("expected data argument")
	} else if len(args) >= 3 {
		aclstr = args[2]
	}

	data := args[1]

	if authUser != "" && authPwd != "" {
		perms, err := client.BuildACL("digest", authUser, authPwd, acls)
		if err != nil {
			return err
		}

		result, err := client.CreateWithACL(path, []byte(data), force, perms)
		if err != nil {
			return err
		}

		log.Infof("Created %+v", result)
	} else {
		result, err := client.Create(path, []byte(data), aclstr, force)
		if err != nil {
			return err
		}

		log.Infof("Created %+v", result)
	}

	return nil
}
