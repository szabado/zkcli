package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsrCmd)
	lsrCmd.PersistentFlags().IntVar(&concurrentRequests, "concurrent_requests", 1, "Number of requests to make in parallel")
}

var lsrCmd = &cobra.Command{
	Use:  "lsr",
	RunE: lsrExecute,
}

func lsrExecute(_ *cobra.Command, _ []string) error {
	children, err := client.ChildrenRecursive(path, concurrentRequests)
	if err != nil {
		return err
	}

	out.PrintStringArray(children)

	return nil
}
