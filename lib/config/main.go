package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config represents the configuration data.
type Config struct {
	Name         string `mapstructure:"name"`
	Environment  string `mapstructure:"environment"`
	Listen       string `mapstructure:"listen"`
	Redis        string `mapstructure:"redis"`
	RedisMaxConn int    `mapstructure:"redis_max_conn"`
	JWTSecret    string `mapstructure:"jwt_secret"`
	SessionTTL   int    `mapstructure:"session_ttl"`
	CookieSecure bool   `mapstructure:"cookie_secure"`
	Key          string `mapstructure:"key"`
	Seed         string `mapstructure:"seed"`
}

// _conf contains defaults used before external config is loaded.
var _conf = Config{
	Listen:       "0.0.0.0:80",
	RedisMaxConn: 30,
	JWTSecret:    "change-me-in-production",
	SessionTTL:   86400,
	CookieSecure: false,
}

// Get returns the global config.
func Get() Config {
	return _conf
}

// Read reads config from an explicit file path or a conventional config directory.
func Read(filePath string) {
	if filePath != "" {
		viper.SetConfigFile(filePath)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetDefault("listen", "0.0.0.0:80")
	viper.SetDefault("redis_max_conn", 30)
	viper.SetDefault("jwt_secret", "change-me-in-production")
	viper.SetDefault("session_ttl", 86400)
	viper.SetDefault("cookie_secure", false)

	_ = viper.BindEnv("name", "APP_NAME")
	_ = viper.BindEnv("environment", "APP_ENVIRONMENT")
	_ = viper.BindEnv("listen", "APP_LISTEN")
	_ = viper.BindEnv("redis", "REDIS_URL")
	_ = viper.BindEnv("redis_max_conn", "REDIS_MAX_CONN")
	_ = viper.BindEnv("jwt_secret", "JWT_SECRET")
	_ = viper.BindEnv("session_ttl", "SESSION_TTL")
	_ = viper.BindEnv("cookie_secure", "COOKIE_SECURE")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config: using defaults because no config file was found: %v", err)
		if err := viper.Unmarshal(&_conf); err != nil {
			log.Printf("config: failed to unmarshal defaults: %v", err)
		}
		return
	}

	if err := viper.Unmarshal(&_conf); err != nil {
		log.Fatal("config: could not unmarshal config\n", err)
	}
}
