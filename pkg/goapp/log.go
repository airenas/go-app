package goapp

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//Log is applications logger
var Log = log.Logger

func initLog() {
	sl := Config.GetString("logger.level")
	if sl == "" {
		sl = "info"
	}
	l, err := zerolog.ParseLevel(strings.ToLower(sl))
	if err != nil {
		Log.Error().Err(err).Msgf("can't parse %s", sl)
	} else {
		Log = log.Logger.Level(l)
	}
	if strings.ToLower(Config.GetString("logger.out")) == "console" {
		Log = log.Logger.Level(l).Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) { w.TimeFormat = "2006-01-02T15:04:05.000" }))
	}
}
