package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymnetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.BindEnv("DB_DRIVER", "DB_DRIVER")
	viper.BindEnv("DB_SOURCE", "DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS", "SERVER_ADDRESS")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY", "TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("ACCESS_TOKEN_DURATION", "ACCESS_TOKEN_DURATION")

	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = nil
		} else {
			return
		}
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}
	return
}
