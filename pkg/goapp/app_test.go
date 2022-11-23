package goapp

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestReadEnvVariable(t *testing.T) {
	os.Setenv("MESSAGESERVER_URL", "olia")
	InitConfig("")
	assert.Equal(t, "olia", Config.GetString("messageServer.url"))
}

func TestReadBoolEnvVariable(t *testing.T) {
	os.Setenv("SENDINFORMMESSAGES", "true")
	InitConfig("")

	assert.Equal(t, true, Config.GetBool("sendInformMessages"))
}

func TestReadConfig(t *testing.T) {
	initAppFromTempFile(t, "messageServer:\n     url: olia\n")

	assert.Equal(t, "olia", Config.GetString("messageServer.url"))
}

func TestEnvBeatsConfig(t *testing.T) {
	os.Setenv("MESSAGESERVER_URL", "xxxx")
	initAppFromTempFile(t, "messageServer:\n     url: olia\n")

	assert.Equal(t, "xxxx", Config.GetString("messageServer.url"))
}

func TestEnvBeatsSubConfig(t *testing.T) {
	os.Setenv("MESSAGESERVER_URL", "xxxx")
	initAppFromTempFile(t, "messageServer:\n     url: olia\n")

	assert.Equal(t, "xxxx", Sub(Config, "messageServer").GetString("url"))
}

func TestSub_NoFail(t *testing.T) {
	initAppFromTempFile(t, "messageServer:\n     url: olia\n")

	assert.NotNil(t, Sub(Config, "messageServer"))
	assert.Nil(t, Sub(Config, "messageServer1"))
}

func TestEnvBeatsSubConfigNested(t *testing.T) {
	os.Setenv("MESSAGE_SERVER_URL", "xxxx")
	initAppFromTempFile(t, "message:\n  server:\n    url: olia\n")

	assert.Equal(t, "xxxx", Sub(Config, "message").GetString("server.url"))
}

// the test fails as where is no option to get current env prefix from config
// func TestEnvBeatsSeveralSubConfigNested(t *testing.T) {
// 	os.Setenv("MESSAGE_SERVER_URL", "xxxx")
// 	initAppFromTempFile(t, "message:\n  server:\n    url: olia\n")

// 	cfg := Sub(Config, "message")

// 	assert.Equal(t, "xxxx", Sub(cfg, "server").GetString("url"))
// }

func TestDefaultLogger(t *testing.T) {
	initDefaultLevel()
	initAppFromTempFile(t, "")

	assert.Equal(t, "info", Log.GetLevel().String())
}

func TestLoggerInitFromConfig(t *testing.T) {
	initDefaultLevel()
	initAppFromTempFile(t, "logger:\n    level: trace\n")

	assert.Equal(t, "trace", Log.GetLevel().String())
}

func TestLoggerLevelInitFromEnv(t *testing.T) {
	initDefaultLevel()

	os.Setenv("LOGGER_LEVEL", "trace")
	initAppFromTempFile(t, "logger:\n    level: info\n")

	assert.Equal(t, "trace", Log.GetLevel().String())
}

func TestStartWitFlags(t *testing.T) {
	f, err := ioutil.TempFile("", "test.*.yml")
	assert.Nil(t, err)
	f.WriteString("logger:\n  level: TRACE")
	f.Sync()
	defer os.Remove(f.Name())

	fs := flag.NewFlagSet("", flag.ExitOnError)
	StartWithFlags(fs, []string{"app", "-c", f.Name()})

	assert.Equal(t, "trace", Log.GetLevel().String())
}

func TestStartWitFlags_Panic(t *testing.T) {
	assert.Panics(t, func() {
		fs := flag.NewFlagSet("", flag.PanicOnError)
		StartWithFlags(fs, []string{"app", "-a", "olia.yml"})
	})
}

func initAppFromTempFile(t *testing.T, data string) {
	f, err := ioutil.TempFile("", "test.*.yml")
	assert.Nil(t, err)
	f.WriteString(data)
	f.Sync()

	defer os.Remove(f.Name())
	InitConfig(f.Name())
}

func initDefaultLevel() {
	Log.Level(zerolog.InfoLevel)
	Log.Output(os.Stdout)
}
