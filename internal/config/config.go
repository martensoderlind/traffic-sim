package config

import (
	"github.com/spf13/viper"
)


type FeatureFlags struct {
    RightOfWaySystem bool `mapstructure:"RIGHT_OF_WAY_SYSTEM"`
}

type Config struct {
    FeatureFlags FeatureFlags `mapstructure:"featureFlags"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigFile("internal/config/config.yaml")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
