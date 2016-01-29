package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var Cfg *viper.Viper

func init() {
	Cfg = viper.New()
}

func LoadConfig() {
	Cfg.SetConfigType("toml")
	Cfg.SetConfigName("config")                  // name of config file (without extension)
	Cfg.AddConfigPath("$HOME/.config/goshortener") // call multiple times to add many search paths
	Cfg.AddConfigPath(".")                       // optionally look for config in the working directory
	Cfg.AutomaticEnv()                           // read in environment variables that match

	if err := Cfg.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
