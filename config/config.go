package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Hostname string `mapstructure:"DB_HOSTNAME"`
	Port     int    `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Database string `mapstructure:"DB_DATABASE"`

	Addr string `mapstructure:"SRV_ADDR"`
}

func Load() *Config {
	env := &Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot read cofiguration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("environment cant be loaded: ", err)
	}

	return env
}
