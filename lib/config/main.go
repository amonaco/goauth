package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config represents the configuration data
type Config struct {
	Name   string `mapstructure:"name"`
	Listen string `mapstructure:"listen"`
	Key    string `mapstructure:"key"`
	Seed   string `mapstructure:"seed"`
}

// Init conf with defaults
var _conf = Config{
	Listen: "0.0.0.0:80",
}

// Get returns the global config
func Get() Config {
	return _conf
}

// Read reads the global config from a json file
func Read(filePath string) {

	// Viper setup
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Bind config to env vars
	viper.BindEnv("name", "APP_NAME")
	viper.BindEnv("listen", "APP_LISTEN")
	viper.BindEnv("postgres", "POSTGRES")

	// Reads the config
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Fatal error config file\n", err)
	}

	err = viper.Unmarshal(&_conf)
	if err != nil {
		log.Fatal("Cound not unmarshall config\n", err)
	}
}
