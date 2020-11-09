package main

import (
	"github.com/airenas/go-app/pkg/goapp"
)

func main() {
	goapp.StartWithDefault()

	goapp.Log.Debug("Now will be printing: Hello world")
	goapp.Log.Info("Hello world!")
}
