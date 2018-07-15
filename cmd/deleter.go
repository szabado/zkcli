package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	deleterCommandUse = "deleter"
)

func init() {
	rootCmd.AddCommand(deleterCmd)
	deleterCmd.PersistentFlags().IntVar(&concurrentRequests, concurrentRequestsFlag, defaultConcurrentRequests, "Number of requests to make in parallel")
}

var deleterCmd = &cobra.Command{
	Use:     deleterCommandUse + " <path>",
	Short:   "Delete the specified znode, as well as any children",
	Aliases: []string{"rmr"},
	RunE:    deleterExecute,
}

func deleterExecute(cmd *cobra.Command, _ []string) error {
	if !force {
		force = true
		log.Warn("%v command requires --force for safety measure", cmd.Use)
	}

	return client.DeleteRecursive(path, concurrentRequests)
}
