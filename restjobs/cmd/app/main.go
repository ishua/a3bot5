package main

import (
	"fmt"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type MyConfig struct {
	Redis string `default:"redis:6379" env:"REDIS" usage:"connect str to redis"`
}

var cfg MyConfig

// init config
func init() {
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"conf/restjobs_config.yaml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println(cfg.Redis)
}
