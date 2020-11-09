package goapp

import (
	"github.com/heirko/go-contrib/logrusHelper"
	"github.com/heralight/logrus_mate"
	"github.com/sirupsen/logrus"
)

//Log is applications logger
var Log = logrus.New()

func initLog() {
	initDefaultLogConfig()
	c := logrusHelper.UnmarshalConfiguration(Config.Sub("logger"))
	initLogFromEnv(&c)
	err := logrusHelper.SetConfig(Log, c)
	if err != nil {
		Log.Error("Can't init log ", err)
	}
}

//initLogFromEnv tries to set level from environment
func initLogFromEnv(c *logrus_mate.LoggerConfig) {
	ll := Config.GetString("logger.level")
	if ll != "" {
		c.Level = ll
	}
}

func initDefaultLogConfig() {
	defaultLogConfig := map[string]interface{}{
		"level":                              "info",
		"formatter.name":                     "text",
		"formatter.options.full_timestamp":   true,
		"formatter.options.timestamp_format": "2006-01-02T15:04:05.000",
	}
	Config.SetDefault("logger", defaultLogConfig)
}
