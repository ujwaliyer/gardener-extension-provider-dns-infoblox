package config

import (
	viper "github.com/spf13/viper"
)

type Config struct {
	Username    string `mapstructure:"IB_TEST_USERNAME"`
	Password    string `mapstructure:"IB_TEST_PASSWORD"`
	Host        string `mapstructure:"IB_TEST_HOST"`
	DefaultZone string `mapstructure:"IB_TEST_DEFAULT_ZONE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
