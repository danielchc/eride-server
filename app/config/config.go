package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server struct {
		Application string `yaml:"application" validate:"required"`
	} `yaml:"server"`
}

var config Configuration

func LoadConfig() (*viper.Viper, error) {
	vp := viper.New()
	vp.SetConfigName("config") // name of config file (without extension)
	vp.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	vp.AddConfigPath(".")
	vp.SetDefault("server.port", 9647)
	vp.SetDefault("security.jwtTokenDuration", 36000000)

	vp.AutomaticEnv()
	vp.SetEnvPrefix("env")
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read error %v", err)
	}
	if err := vp.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshall the config %v", err)
	}

	// validate := validator.New()
	// if err := validate.Struct(&config); err != nil {
	// 	return nil, fmt.Errorf("missing required attributes %v\n", err)
	// }

	return vp, nil
}
