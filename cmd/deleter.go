package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
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
		force = true
		log.Warn("%v command requires --force for safety measure", cmd.Use)
	}

	return client.DeleteRecursive(path, concurrentRequests)
}
