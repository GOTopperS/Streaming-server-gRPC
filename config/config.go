package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AwsAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	Port               int    `mapstructure:"PORT"`
}

func LoadConfig(path string) (*Config, error) {

	for k, v := range Defaults() {
		viper.SetDefault(k, v)
	}

	viper.AddConfigPath(path)
	viper.SetConfigFile("app.env")

	/*
		Tells viper to load values from the environment variables
		and overwrite the ones loaded from the file
	*/
	viper.AutomaticEnv()

	// tells viper to start reading the config
	// if err = viper.ReadInConfig(); err != nil {
	// 	return
	// }
	config := &Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Defaults() map[string]string {
	return map[string]string{
		"AWS_ACCESS_KEY_ID":     "",
		"AWS_SECRET_ACCESS_KEY": "",
		"PORT":                  "8080",
	}
}
