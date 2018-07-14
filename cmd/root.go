package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fJancsoSzabo/zkcli/output"
	"github.com/fJancsoSzabo/zkcli/zk"
)

const (
	txtFormat  = "txt"
	jsonFormat = "json"

	zkcliCommandUse = "zkcli"

	serverFlag             = "servers"
	omitNewlineFlag        = "n"
	verboseFlag            = "verbose"
	debugFlag              = "debug"
	concurrentRequestsFlag = "concurrent_requests"
	formatFlag             = "format"
	forceFlag              = "force"
	authUserFlag           = "auth_usr"
	authPwdFlag            = "auth_pwd"

	defaultConcurrentRequests = 1
	defaultFormat             = txtFormat
	defaultAuthUser           = ""
	defaultAuthPwd            = ""
	defaultDebug              = false
	defaultVerbose            = false
	defaultOmitnewline        = false
	defaultPath               = ""
	defaultForce              = false
	defaultServer             = ""
)

const (
	aclRead = 1 << iota
	aclWrite
	aclCreate
	aclDelete
	aclAdmin
	aclAll = aclRead | aclWrite | aclCreate | aclDelete | aclAdmin
)

var (
	// Flag variables
	servers            string
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
	stdin  io.Reader
	osExit func(code int)
)

func init() {
	stdin = os.Stdin
	osExit = os.Exit

	rootCmd.PersistentFlags().String(serverFlag, defaultServer, "srv1[:port1][,srv2[:port2]...]")
	rootCmd.PersistentFlags().BoolVar(&force, forceFlag, defaultForce, "force operation")
	rootCmd.PersistentFlags().StringVar(&format, formatFlag, defaultFormat, "output format ("+txtFormat+"|"+jsonFormat+")")
	rootCmd.PersistentFlags().BoolVar(&omitNewline, omitNewlineFlag, defaultOmitnewline, "omit trailing newline")
	rootCmd.PersistentFlags().BoolVar(&verbose, verboseFlag, defaultVerbose, "verbose")
	rootCmd.PersistentFlags().BoolVar(&debug, debugFlag, defaultDebug, "debug mode (very verbose)")
	rootCmd.PersistentFlags().String(authUserFlag, defaultAuthUser, "optional, digest scheme, user")
	rootCmd.PersistentFlags().String(authPwdFlag, defaultAuthPwd, "optional, digest scheme, pwd")

	viper.BindPFlag(serverFlag, rootCmd.PersistentFlags().Lookup(serverFlag))
	viper.BindPFlag(authUserFlag, rootCmd.PersistentFlags().Lookup(authUserFlag))
	viper.BindPFlag(authPwdFlag, rootCmd.PersistentFlags().Lookup(authPwdFlag))
	viper.BindEnv(serverFlag, "ZKCLI_SERVERS")
	viper.BindEnv(authUserFlag, "ZKCLI_AUTH_USER")
	viper.BindEnv(authPwdFlag, "ZKCLI_AUTH_PWD")
	servers = viper.Get(serverFlag).(string)
	authUser = viper.Get(authUserFlag).(string)
	authPwd = viper.Get(authPwdFlag).(string)
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

		serversArray := strings.Split(servers, ",")
		if len(serversArray) == 0 || serversArray[0] == "" {
			return errors.Errorf("Expected comma delimited list of servers via --servers")
		}

		if len(args) == 0 {
			return errors.Errorf("Path must be specified")
		}
		path = args[0]

		if strings.HasSuffix(path, "/") && path != "/" {
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

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		osExit(1)
	}
}
