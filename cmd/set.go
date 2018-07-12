package cmd

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	setCommandUse = "set"
)

var (
	data []byte
)

func init() {
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use: setCommandUse,
	PreRunE: func(_ *cobra.Command, args []string) error {
		if len(args) >= 2 {
			data = []byte(args[1])
		} else {
			var err error
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
		}
		return nil
	},
	RunE: setExecute,
}

func setExecute(_ *cobra.Command, _ []string) error {
	result, err := client.Set(path, data)
	if err != nil {
		return err
	}

	log.Infof("Set %+v", result)

	return nil
}
