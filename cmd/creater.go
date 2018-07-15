package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	createrCommandUse = "creater"
)

func init() {
	rootCmd.AddCommand(createrCmd)
	createrCmd.PersistentFlags().StringVar(&acls, aclsFlag, fmt.Sprint(defaultAcls), fmt.Sprintf("optional, csv list [%v|,%v|,%v|,%v|,%v|,%v]", aclRead, aclWrite, aclCreate, aclDelete, aclAdmin, aclAll))
}

var createrCmd = &cobra.Command{
	Use:   createrCommandUse + " <path> [acls]",
	Short: "Create the specified znode, as well as any required parents",
	RunE:  createrExecute,
}

func createrExecute(cmd *cobra.Command, args []string) error {
	force = true
	return createExecute(cmd, args)
}
