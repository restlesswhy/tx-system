package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Hostname string `mapstructure:"DB_HOSTNAME"`
	Port     int    `mapstructure:"DB_PORT"`
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	Database string `mapstructure:"POSTGRES_DB"`

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
