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

	const (
		testPath = "/test"
		testData = "data"
	)
	hostsArg := strings.Join(hosts, ",")

	os.Args = []string{zkcliCommandUse, createCommandUse, testPath, testData, "--" + serverFlag, hostsArg}
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
	os.Args = []string{zkcliCommandUse, getCommandUse, testPath, "--" + serverFlag, hostsArg}
	err = rootCmd.Execute()
	require.NoError(err)
	reader := bufio.NewReader(r)
	output, _ := reader.ReadString('\n')
	require.NoError(err)
	assert.Equal(testData + "\n", output)

	//os.Args = []string{zkcliCommandUse, getCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag}
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

	os.Args = []string{zkcliCommandUse, setCommandUse, testPath, updatedTestData, "--" + serverFlag, hostsArg}
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte(updatedTestData), value)

	os.Args = []string{zkcliCommandUse, setCommandUse, testPath, "", "--" + serverFlag, hostsArg}
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte{}, value)

	os.Args = []string{zkcliCommandUse, deleteCommandUse, testPath, "", "--" + serverFlag, hostsArg}
	err = rootCmd.Execute()
	require.NoError(err)

	exists, stat, err := zkConn.Exists(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.False(exists)
}

func TestCRUDRecurisve(t *testing.T) {
	require := r.New(t)
	assert := a.New(t)

	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()

	const (
		baseTestPath = "/test"
		testPath = baseTestPath + "/example/debugging?"
		testData = "data"
	)
	hostsArg := strings.Join(hosts, ",")

	os.Args = []string{zkcliCommandUse, createrCommandUse, testPath, testData, "--" + serverFlag, hostsArg}
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
	os.Args = []string{zkcliCommandUse, getCommandUse, testPath, "--" + serverFlag, hostsArg}
	err = rootCmd.Execute()
	require.NoError(err)
	reader := bufio.NewReader(r)
	output, _ := reader.ReadString('\n')
	require.NoError(err)
	assert.Equal(testData + "\n", output)

	//os.Args = []string{zkcliCommandUse, getCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag}
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

	os.Args = []string{zkcliCommandUse, setCommandUse, testPath, updatedTestData, "--" + serverFlag, hostsArg}
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte(updatedTestData), value)

	os.Args = []string{zkcliCommandUse, setCommandUse, testPath, "", "--" + serverFlag, hostsArg}
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal([]byte{}, value)

	os.Args = []string{zkcliCommandUse, deleterCommandUse, baseTestPath, "", "--" + serverFlag, hostsArg}
	err = rootCmd.Execute()
	require.NoError(err)

	exists, stat, err := zkConn.Exists(baseTestPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.False(exists)
}

func TestCreate(t *testing.T) {
	require := r.New(t)
	//assert := a.New(t)

	hosts, id, err := StartServer()
	require.NoError(err)
	defer id.KillRemove()
	zkConn, _, err := zookeeper.Connect(hosts, time.Hour)
	defer zkConn.Close()

	client = zk.NewZooKeeper()
	client.SetServers(hosts)

	const (
		testPath = "/test"
		//testData = "data"
	)

	os.Args = []string{zkcliCommandUse, createCommandUse, testPath}
	err = rootCmd.Execute()
	require.Error(err)

	os.Args = []string{zkcliCommandUse, createCommandUse, "invalidPath"}
	err = rootCmd.Execute()
	require.Error(err)
}