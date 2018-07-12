package cmd

import (
	"github.com/ory/dockertest"
	"time"
	"github.com/sirupsen/logrus"
	"testing"
	r "github.com/stretchr/testify/require"
	a "github.com/stretchr/testify/assert"
	zookeeper "github.com/samuel/go-zookeeper/zk"
	"github.com/fJancsoSzabo/zkcli/zk"
	"os"
	"strings"
	"bufio"
)

const (
	ServerPollingInterval = 10 * time.Millisecond
)

type logger struct {
}

func (l *logger) Printf(message string, values ...interface{}) {
	logrus.StandardLogger().Infof(message, values)
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
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	reader := bufio.NewReader(r)
	output, _ := reader.ReadString('\n')
	require.NoError(err)
	assert.Equal(testData + "\n", output)

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

	rootCmd.SetArgs([]string{setCommandUse, testPath, updatedTestData, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte(updatedTestData), value)

	rootCmd.SetArgs([]string{setCommandUse, testPath, "", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte{}, value)

	rootCmd.SetArgs([]string{deleteCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	exists, stat, err := zkConn.Exists(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.False(exists)

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
		testPath = baseTestPath + "/example/debugging?"
		testData = "data"
	)

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
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	reader := bufio.NewReader(r)
	output, _ := reader.ReadString('\n')
	require.NoError(err)
	assert.Equal(testData + "\n", output)

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

	rootCmd.SetArgs([]string{setCommandUse, testPath, updatedTestData, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte(updatedTestData), value)

	rootCmd.SetArgs([]string{setCommandUse, testPath, "", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte{}, value)

	rootCmd.SetArgs([]string{deleterCommandUse, baseTestPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	exists, stat, err := zkConn.Exists(baseTestPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.False(exists)

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
		//testData = "data"
	)

	rootCmd.SetArgs([]string{createCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)

	rootCmd.SetArgs([]string{createCommandUse, "/../invalidPath", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)

	exists, stat, err := zkConn.Exists("/invalidPath")
	require.Nil(err)
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

	rootCmd.SetArgs([]string{setCommandUse, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)

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
	rootCmd.SetArgs([]string{createCommandUse, "/path", "data"})
	err = rootCmd.Execute()
	require.Error(err)

	exists, stat, err := zkConn.Exists("/path")
	require.Nil(err)
	assert.NotNil(stat)
	assert.False(exists)

	// Test verbose flag
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

	rootCmd.SetArgs([]string{existsCommandUse, "/../invalidPath", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)
}
