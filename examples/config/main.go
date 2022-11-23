package main

import (
	"github.com/airenas/go-app/pkg/goapp"
)

func main() {
	goapp.StartWithDefault()

	name := goapp.Config.GetString("sample_name")
	goapp.Log.Info().Str("name", name).Msg("Hello world!")
}
