package goapp

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// InitConfig tries to load config.yml from exe's dir
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
			return fmt.Errorf("can't get the app directory: %w", err)
		}
		Config.AddConfigPath(filepath.Dir(ex))
		Config.SetConfigName("config")
	}

	if err := Config.ReadInConfig(); err != nil {
		Log.Warn().Err(err).Msg("can't read config")
		if failOnNoFail {
			return fmt.Errorf("can't read config: %w", err)
		}
	}
	initLog()
	if Config.ConfigFileUsed() != "" {
		Log.Info().Str("file", Config.ConfigFileUsed()).Msg("Config loaded")
	}
	return nil
}

// StartWithDefault default app initialization function
// Tries to load config from commandline option '-c'
// panics on error
func StartWithDefault() {
	StartWithFlags(flag.CommandLine, os.Args)
}

// StartWithFlags app initialization function with flagset
// panics on error
func StartWithFlags(fs *flag.FlagSet, args []string) {
	cFile := fs.String("c", "", "Config yml file")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:[params] \n", args[0])
		fs.PrintDefaults()
	}
	if err := fs.Parse(args[1:]); err != nil {
		Log.Fatal().Err(err).Msg("can't init app")
	}

	if err := InitConfig(*cFile); err != nil {
		Log.Fatal().Err(err).Msg("can't init app")
	}
}
