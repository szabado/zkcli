package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fJancsoSzabo/zkcli/output"
	"github.com/fJancsoSzabo/zkcli/zk"
)

const (
	txtFormat  = "txt"
	jsonFormat = "json"

	zkcliCommandUse = "zkcli"

	serverFlag = "servers"
	omitNewlineFlag = "n"
)

var (
	// Flag variables
	servers            string
	command            string
	force              bool
	format             string
	omitNewline        bool
	verbose            bool
	debug              bool
	authUser           string
	authPwd            string
	concurrentRequests int
	path               string

	client *zk.ZooKeeper
	out    output.Printer
)

func init() {
	rootCmd.PersistentFlags().StringVar(&servers, serverFlag, "", "srv1[:port1][,srv2[:port2]...]")
	rootCmd.PersistentFlags().BoolVar(&force, "force", false, "force operation")
	rootCmd.PersistentFlags().StringVar(&format, "format", txtFormat, "output format ("+txtFormat+"|"+jsonFormat+")")
	rootCmd.PersistentFlags().BoolVar(&omitNewline, omitNewlineFlag, false, "omit trailing newline")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode (very verbose)")
	rootCmd.PersistentFlags().StringVar(&authUser, "auth_usr", "", "optional, digest scheme, user")
	rootCmd.PersistentFlags().StringVar(&authPwd, "auth_pwd", "", "optional, digest scheme, pwd")

}

var rootCmd = &cobra.Command{
	Use:   zkcliCommandUse,
	Short: "A CLI to interact with Zookeeper",
	PersistentPreRunE: func(_ *cobra.Command, args []string) error {
		log.SetLevel(log.ErrorLevel)
		if verbose {
			log.SetLevel(log.InfoLevel)
		}
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		switch format {
		case txtFormat:
			out = &output.TxtPrinter{
				OmitTrailingNL: omitNewline,
			}
		case jsonFormat:
			out = &output.JSONPrinter{}
		default:
			return errors.Errorf("unknown output type %s", format)
		}

		log.Info("starting")

		serversArray := strings.Split(servers, ",")
		if len(serversArray) == 0 {
			log.Fatal("Expected comma delimited list of servers via --servers")
		}

		if strings.HasSuffix(path, "/") {
			log.Warn("Removing trailing / from path")
			path = strings.TrimSuffix(path, "/")
		}

		rand.Seed(time.Now().UnixNano())
		client = zk.NewZooKeeper()
		client.SetServers(serversArray)

		if authUser != "" && authPwd != "" {
			authExp := fmt.Sprintf("%v:%v", authUser, authPwd)
			client.SetAuth("digest", []byte(authExp))
		}

		path = args[0]

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
