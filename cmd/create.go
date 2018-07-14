package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	createCommandUse = "create"

	aclsFlag      = "acls"
	defaultAclstr = ""
	defaultAcls   = aclAll
)

var (
	aclstr string
	acls   string
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVar(&acls, aclsFlag, fmt.Sprint(defaultAcls), fmt.Sprintf("optional, csv list [%v|,%v|,%v|,%v|,%v|,%v]", aclRead, aclWrite, aclCreate, aclDelete, aclAdmin, aclAll))
}

var createCmd = &cobra.Command{
	Use:   createCommandUse,
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
		if aclstr != "" {
			log.Warnf("Ignoring acl %v since --%v and --%v were specified", aclstr, authUserFlag, authPwdFlag)
		}

		perms, err := client.BuildACL("digest", authUser, authPwd, acls)
		if err != nil {
			return err
		}

		result, err := client.CreateWithACL(path, []byte(data), force, perms)
		if err != nil {
			return err
		}

		out.Printf("Created %+v", result)
	} else {
		if authUser != "" || authPwd != "" {
			log.Warnf("Both --%v and --%v must be specified", authUserFlag, authPwdFlag)
		}

		result, err := client.Create(path, []byte(data), aclstr, force)
		if err != nil {
			return err
		}

		out.Printf("Created %+v", result)
	}

	return nil
}
