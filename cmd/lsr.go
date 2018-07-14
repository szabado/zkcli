package cmd

import (
	"github.com/spf13/cobra"
)

const(
	lsrCommandUse = "lsr"
)

func init() {
	rootCmd.AddCommand(lsrCmd)
	lsrCmd.PersistentFlags().IntVar(&concurrentRequests, concurrentRequestsFlag, defaultConcurrentRequests, "Number of requests to make in parallel")
}

var lsrCmd = &cobra.Command{
	Use:   lsrCommandUse,
	Short: "Print the children of the current znode recursively",
	RunE:  lsrExecute,
}

func lsrExecute(_ *cobra.Command, _ []string) error {
	children, err := client.ChildrenRecursive(path, concurrentRequests)
	if err != nil {
		return err
	}

	out.PrintArray(children)

	return nil
}
