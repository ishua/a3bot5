package main

import (
	"fmt"

	"github.com/cristalhq/aconfig"
)

type MyConfig struct {
	Token string `required:"true" env:"TELEGRAMBOTTOKEN" usage:"token for your telegram bot"`
	Debug bool   `default:"false" usage:"turn on debug mode"`
}

var cfg MyConfig

func main() {
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"config.json"},
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}

	fmt.Println(cfg.Token)
}
