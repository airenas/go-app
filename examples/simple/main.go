package main

import (
	"fmt"

	"github.com/airenas/go-app/pkg/goapp"
)

func main() {
	goapp.StartWithDefault()

	goapp.Log.Debug().Msg("Now will be printing: Hello world")
	goapp.Log.Info().Msg("Hello world!")
	goapp.Log.Error().Err(fmt.Errorf("Ops")).Send()
}
