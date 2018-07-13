package cmd

import (
	"github.com/spf13/cobra"
)

const (
	createrCommandUse = "creater"
)

func init() {
	rootCmd.AddCommand(createrCmd)
	createrCmd.PersistentFlags().StringVar(&acls, aclsFlag, defaultAcls, "optional, csv list [1|,2|,4|,8|,16|,31]")
}

var createrCmd = &cobra.Command{
	Use:  createrCommandUse,
	Short: "Create the specified znode, as well as any required parents",
	RunE: createrExecute,
}

func createrExecute(cmd *cobra.Command, args []string) error {
	force = true
	return createExecute(cmd, args)
}
