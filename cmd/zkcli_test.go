package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest"
	zookeeper "github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
	a "github.com/stretchr/testify/assert"
	r "github.com/stretchr/testify/require"

	"github.com/fJancsoSzabo/zkcli/zk"
)

const (
	ServerPollingInterval = 10 * time.Millisecond
)

type logger struct {
}

func (l *logger) Printf(message string, values ...interface{}) {
	logrus.StandardLogger().Infof(message, values)
}

func loadDefaultValues() {
	aclstr = defaultAclstr
	acls = fmt.Sprint(defaultAcls)
	servers = defaultServer
	force = defaultForce
	format = defaultFormat
	omitNewline = defaultOmitnewline
	verbose = defaultVerbose
	debug = defaultDebug
	authUser = defaultAuthUser
	authPwd = defaultAuthPwd
	concurrentRequests = defaultConcurrentRequests
	path = defaultPath

	client = nil
	out = nil
}

func StartServer() (hosts []string, id dockertest.ContainerID, err error) {
	id, err = dockertest.ConnectToZooKeeper(10, ServerPollingInterval, func(url string) bool {
		hosts = []string{url}
		conn, _, err := zookeeper.Connect([]string{url}, time.Second, zookeeper.WithLogger(&logger{}))
		if err != nil {
			return false
		}
		conn.Close()

		return true
	})

	return hosts, id, err
}

func TestCRUD(t *testing.T) {
	require := r.New(t)
	assert := a.New(t)

	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()
	hostsArg := strings.Join(hosts, ",")

	const (
		testPath = "/test"
		testData = "data"
	)

	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, testPath, testData, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err := zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte(testData), value)

	tempOutput := os.Stdout
	r, w, err := os.Pipe()
	require.Nil(err)
	defer r.Close()
	defer w.Close()
	os.Stdout = w
	defer func() {
		os.Stdout = tempOutput
	}()

	loadDefaultValues()
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	reader := bufio.NewReader(r)
	output, _ := reader.ReadString('\n')
	require.NoError(err)
	assert.Equal(testData+"\n", output)

	//rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag})
	//err = rootCmd.Execute()
	//require.NoError(err)
	//os.Stdout = tempOutput
	//outputBytes := make([]byte, 100)
	//bytesRead, err := reader.Read(outputBytes)
	//require.Equal(len(testData), bytesRead)
	//require.NoError(err)
	//assert.Equal(testData, string(outputBytes[:bytesRead]))

	const (
		updatedTestData = "newVal"
	)

	loadDefaultValues()
	rootCmd.SetArgs([]string{setCommandUse, testPath, updatedTestData, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte(updatedTestData), value)

	loadDefaultValues()
	rootCmd.SetArgs([]string{setCommandUse, testPath, "", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte{}, value)

	loadDefaultValues()
	rootCmd.SetArgs([]string{deleteCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	exists, stat, err := zkConn.Exists(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.False(exists)

	loadDefaultValues()
	rootCmd.SetArgs([]string{existsCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
}

func TestCRUDRecurisve(t *testing.T) {
	require := r.New(t)
	assert := a.New(t)

	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()
	hostsArg := strings.Join(hosts, ",")

	const (
		baseTestPath = "/test"
		testPath     = baseTestPath + "/example/debugging?"
		testData     = "data"
	)

	loadDefaultValues()
	rootCmd.SetArgs([]string{createrCommandUse, testPath, testData, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err := zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte(testData), value)

	tempOutput := os.Stdout
	r, w, err := os.Pipe()
	require.Nil(err)
	defer r.Close()
	defer w.Close()
	os.Stdout = w
	defer func() {
		os.Stdout = tempOutput
	}()

	loadDefaultValues()
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	reader := bufio.NewReader(r)
	output, _ := reader.ReadString('\n')
	require.NoError(err)
	assert.Equal(testData+"\n", output)

	//rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag})
	//err = rootCmd.Execute()
	//require.NoError(err)
	//os.Stdout = tempOutput
	//outputBytes := make([]byte, 100)
	//bytesRead, err := reader.Read(outputBytes)
	//require.Equal(len(testData), bytesRead)
	//require.NoError(err)
	//assert.Equal(testData, string(outputBytes[:bytesRead]))

	const (
		updatedTestData = "newVal"
	)

	loadDefaultValues()
	rootCmd.SetArgs([]string{setCommandUse, testPath, updatedTestData, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte(updatedTestData), value)

	loadDefaultValues()
	rootCmd.SetArgs([]string{setCommandUse, testPath, "", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte{}, value)

	loadDefaultValues()
	rootCmd.SetArgs([]string{deleterCommandUse, baseTestPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	exists, stat, err := zkConn.Exists(baseTestPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.False(exists)

	loadDefaultValues()
	rootCmd.SetArgs([]string{existsCommandUse, baseTestPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
}

func TestCreate(t *testing.T) {
	require := r.New(t)
	assert := a.New(t)

	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()
	hostsArg := strings.Join(hosts, ",")

	client = zk.NewZooKeeper()
	client.SetServers(hosts)

	const (
		testPath = "/test"
		testData = "data"
	)

	// No data provided
	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)

	// Invalid path
	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, "/../invalidPath", testData, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)

	exists, stat, err := zkConn.Exists("/invalidPath")
	require.Nil(err)
	assert.NotNil(stat)
	assert.False(exists)

	// Create with credentials
	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, "/authpath", testData, "--" + serverFlag, hostsArg, "--" + authUserFlag, "testuser", "--" + authPwdFlag, "testpassword"})
	err = rootCmd.Execute()
	require.NoError(err)

	exists, stat, err = zkConn.Exists("/authpath")
	require.NoError(err)
	assert.NotNil(stat)
	assert.True(exists)

	value, stat, err := zkConn.Get("/authpath")
	require.Error(err)
	assert.Equal(&zookeeper.Stat{}, stat)
	assert.NotEqual(testData, string(value))

	// Try to create with credentials and an invalid path
	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, "/../invalidPath", testData, "--" + serverFlag, hostsArg, "--" + authUserFlag, "testuser", "--" + authPwdFlag, "testpassword"})
	err = rootCmd.Execute()
	require.Error(err)

	// Try to create with credentials and an invalid acl path
	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, "/authpaththesecond", "--" + aclsFlag, "What is your favourite colour?", testData, "--" + serverFlag, hostsArg, "--" + authUserFlag, "testuser", "--" + authPwdFlag, "testpassword"})
	err = rootCmd.Execute()
	require.Error(err)

	exists, stat, err = zkConn.Exists("/authpaththesecond")
	require.NoError(err)
	assert.NotNil(stat)
	assert.False(exists)
}

func TestSet(t *testing.T) {
	require := r.New(t)

	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()
	hostsArg := strings.Join(hosts, ",")

	client = zk.NewZooKeeper()
	client.SetServers(hosts)

	loadDefaultValues()
	rootCmd.SetArgs([]string{setCommandUse, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)

	loadDefaultValues()
	rootCmd.SetArgs([]string{setCommandUse, "/../invalidPath", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)
}

func TestRoot(t *testing.T) {
	require := r.New(t)
	assert := a.New(t)
	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()
	hostsArg := strings.Join(hosts, ",")

	client = zk.NewZooKeeper()
	client.SetServers(hosts)
	rootCmd.SilenceUsage = true

	// No server specified
	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, "/path", "data"})
	err = rootCmd.Execute()
	require.Error(err)

	exists, stat, err := zkConn.Exists("/path")
	require.Nil(err)
	assert.NotNil(stat)
	assert.False(exists)

	// Test verbose flag
	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, "/path", "--" + verboseFlag, "--" + serverFlag, hostsArg, "data"})
	err = rootCmd.Execute()
	require.NoError(err)
	assert.True(verbose)
	assert.Equal(logrus.InfoLevel, logrus.GetLevel())

	value, stat, err := zkConn.Get("/path")
	require.Nil(err)
	assert.NotNil(stat)
	assert.Equal("data", string(value))

	// debug flag and an extra slash at the end of the path
	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, "/path/nested/", "--" + debugFlag, "--" + serverFlag, hostsArg, "even more data!"})
	err = rootCmd.Execute()
	require.NoError(err)
	assert.True(debug)
	assert.Equal(logrus.DebugLevel, logrus.GetLevel())

	value, stat, err = zkConn.Get("/path/nested")
	require.Nil(err)
	assert.NotNil(stat)
	assert.Equal("even more data!", string(value))

}

func TestGet(t *testing.T) {
	require := r.New(t)

	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()
	hostsArg := strings.Join(hosts, ",")

	client = zk.NewZooKeeper()
	client.SetServers(hosts)

	loadDefaultValues()
	rootCmd.SetArgs([]string{getCommandUse, "/../invalidPath", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)
}

func TestExists(t *testing.T) {
	require := r.New(t)

	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()
	hostsArg := strings.Join(hosts, ",")

	client = zk.NewZooKeeper()
	client.SetServers(hosts)

	loadDefaultValues()
	rootCmd.SetArgs([]string{existsCommandUse, "/../invalidPath", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)
}
