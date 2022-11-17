package config

import (
	viper "github.com/spf13/viper"
)

type Config struct {
	Username    string `mapstructure:"INFOBLOX_USERNAME"`
	Password    string `mapstructure:"INFOBLOX_PASSWORD"`
	Host        string `mapstructure:"INFOBLOX_HOST"`
	DefaultZone string `mapstructure:"IB_TEST_DEFAULT_ZONE"`
}

var cfg Config

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("../../../../config")
	viper.SetConfigName("infoblox")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	cfg = config
	return
}

func GetConfig() (config Config) {
	return cfg
}
