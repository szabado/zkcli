package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest"
	zookeeper "github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
	a "github.com/stretchr/testify/assert"
	r "github.com/stretchr/testify/require"

	"github.com/fJancsoSzabo/zkcli/output"
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

func loadDefaultValues() (stdoutBuf *bytes.Buffer, stdinBuf *bytes.Buffer) {
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
	stdoutBuf = new(bytes.Buffer)
	output.Out = stdoutBuf

	stdinBuf = new(bytes.Buffer)
	stdin = stdinBuf
	return stdoutBuf, stdinBuf
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

	output, _ := loadDefaultValues()
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err := output.ReadString('\n')
	require.NoError(err)
	assert.Equal(testData+"\n", val)

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag})
	err = rootCmd.Execute()
	require.NoError(err)

	val, err = output.ReadString('\n')
	require.Equal(io.EOF, err)
	assert.Equal(testData, val)

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

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{existsCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = output.ReadString('\n')
	require.NoError(err)
	require.Equal(fmt.Sprintln(false), val)
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

	output, _ := loadDefaultValues()
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err := output.ReadString('\n')
	require.NoError(err)
	assert.Equal(testData+"\n", val)

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = output.ReadString('\n')
	require.Equal(io.EOF, err)
	assert.Equal(testData, val)

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

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{existsCommandUse, baseTestPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = output.ReadString('\n')
	require.NoError(err)
	require.Equal(fmt.Sprintln(false), val)

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
	assert := a.New(t)

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

	rootCmd.SetArgs([]string{createCommandUse, "/path", "default testing value", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	_, in := loadDefaultValues()
	in.WriteString("data")
	rootCmd.SetArgs([]string{setCommandUse, "/path", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err := zkConn.Get("/path")
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal("data", string(value))
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

func TestLs(t *testing.T) {
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
	rootCmd.SetArgs([]string{createrCommandUse, "/p/a/t/h/s", "test value one", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	loadDefaultValues()
	rootCmd.SetArgs([]string{createrCommandUse, "/p/otoooooooo", "--" + serverFlag, hostsArg, "test value two"})
	err = rootCmd.Execute()
	require.NoError(err)

	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, "/xyz", "test value three", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	output, _ := loadDefaultValues()
	rootCmd.SetArgs([]string{lsCommandUse, "/", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err := ioutil.ReadAll(output)
	require.NoError(err)
	require.Equal("p\nxyz\nzookeeper\n", string(val))

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{lsCommandUse, "/p", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = ioutil.ReadAll(output)
	require.NoError(err)
	require.Equal("a\notoooooooo\n", string(val))

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{lsCommandUse, "/p/a", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = ioutil.ReadAll(output)
	require.NoError(err)
	require.Equal("t\n", string(val))

	loadDefaultValues()
	rootCmd.SetArgs([]string{lsCommandUse, "/../invalidpath/", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.Error(err)

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{lsrCommandUse, "/", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = ioutil.ReadAll(output)
	require.NoError(err)
	require.Equal("p\np/a\np/a/t\np/a/t/h\np/a/t/h/s\np/otoooooooo\nxyz\nzookeeper\nzookeeper/quota\n", string(val))

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{lsrCommandUse, "/p", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = ioutil.ReadAll(output)
	require.NoError(err)
	require.Equal("a\na/t\na/t/h\na/t/h/s\notoooooooo\n", string(val))

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{lsrCommandUse, "/p/a", "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = ioutil.ReadAll(output)
	require.NoError(err)
	require.Equal("t\nt/h\nt/h/s\n", string(val))

	loadDefaultValues()
	rootCmd.SetArgs([]string{lsrCommandUse, "/../invalidpath/", "--" + serverFlag, hostsArg, "--" + debugFlag})
	err = rootCmd.Execute()
	require.Error(err)
}

func TestAcls(t *testing.T) {
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
		testData = "pigeon"
		acls1    = "world:anyone:rwa"
		acls2    = "world:anyone:rwa,digest:someuser:hashedpw:cdrwa"
	)

	loadDefaultValues()
	rootCmd.SetArgs([]string{createCommandUse, testPath, testData, acls1, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err := zkConn.Get("/test")
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal(testData, string(value))

	output, _ := loadDefaultValues()
	rootCmd.SetArgs([]string{getCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err := ioutil.ReadAll(output)
	require.NoError(err)
	assert.Equal(testData, string(val))

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{getAclCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = ioutil.ReadAll(output)
	require.NoError(err)
	assert.Equal(strings.Replace(acls1, ",", "\n", -1), string(val))

	loadDefaultValues()
	rootCmd.SetArgs([]string{getAclCommandUse, "/../invalidPath", "--" + serverFlag, hostsArg, "--" + omitNewlineFlag})
	err = rootCmd.Execute()
	require.Error(err)

	loadDefaultValues()
	rootCmd.SetArgs([]string{setAclCommandUse, "/../invalidPath", acls1, "--" + serverFlag, hostsArg, "--" + omitNewlineFlag})
	err = rootCmd.Execute()
	require.Error(err)

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{setAclCommandUse, testPath, acls2, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)

	value, stat, err = zkConn.Get(testPath)
	require.NoError(err)
	assert.NotNil(stat)
	assert.Equal(testData, string(value))

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{getAclCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = ioutil.ReadAll(output)
	require.NoError(err)
	assert.Equal(strings.Replace(acls2+"\n", ",", "\n", -1), string(val))

	_, input := loadDefaultValues()
	input.WriteString(acls1)
	rootCmd.SetArgs([]string{setAclCommandUse, testPath, "--" + serverFlag, hostsArg, "--" + debugFlag})
	err = rootCmd.Execute()
	require.NoError(err)

	output, _ = loadDefaultValues()
	rootCmd.SetArgs([]string{getAclCommandUse, testPath, "--" + serverFlag, hostsArg})
	err = rootCmd.Execute()
	require.NoError(err)
	val, err = ioutil.ReadAll(output)
	require.NoError(err)
	assert.Equal(acls1+"\n", string(val))
}
