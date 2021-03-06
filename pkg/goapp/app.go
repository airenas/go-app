package goapp

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

//InitConfig tries to load config.yml from exe's dir
func InitConfig(configFile string) error {
	InitEnv(Config)

	failOnNoFail := false
	if configFile != "" {
		// Use config file from the flag.
		Config.SetConfigFile(configFile)
		failOnNoFail = true
	} else {
		// Find home directory.
		ex, err := os.Executable()
		if err != nil {
			return errors.Wrap(err, "Can't get the app directory")
		}
		Config.AddConfigPath(filepath.Dir(ex))
		Config.SetConfigName("config")
	}

	if err := Config.ReadInConfig(); err != nil {
		Log.Warn("Can't read config:", err)
		if failOnNoFail {
			return errors.Wrap(err, "Can't read config:")
		}
	}
	initLog()
	if Config.ConfigFileUsed() != "" {
		Log.Info("Config loaded from: ", Config.ConfigFileUsed())
	}
	return nil
}

//StartWithDefault default app initialization function
// Tries to load config from commandline option '-c'
// panics on error
func StartWithDefault() {
	StartWithFlags(flag.CommandLine, os.Args)
}

//StartWithFlags app initialization function with flagset
// panics on error
func StartWithFlags(fs *flag.FlagSet, args []string) {
	cFile := fs.String("c", "", "Config yml file")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:[params] \n", args[0])
		fs.PrintDefaults()
	}
	fs.Parse(args[1:])
	err := InitConfig(*cFile)
	if err != nil {
		Log.Fatal(errors.Wrap(err, "Can't init app"))
	}
}
