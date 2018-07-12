package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	deleterCommandUse = "deleter"
)

func init() {
	rootCmd.AddCommand(deleterCmd)
	deleterCmd.PersistentFlags().IntVar(&concurrentRequests, "concurrent_requests", 1, "Number of requests to make in parallel")
}

var deleterCmd = &cobra.Command{
	Use:  deleterCommandUse,
	Short: "Delete the specified znode, as well as any children",
	Aliases: []string{"rmr"},
	RunE: deleterExecute,
}

func deleterExecute(cmd *cobra.Command, _ []string) error {
	if !force {
		log.Fatal(cmd.Use + " command requires --force for safety measure")
	}

	return client.DeleteRecursive(path, concurrentRequests)
}
