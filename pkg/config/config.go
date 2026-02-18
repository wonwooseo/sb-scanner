package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func ReadConfig(cfgF string) (*viper.Viper, error) {
	v := viper.New()
	if cfgF != "" {
		v.SetConfigFile(cfgF)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	} else { // deployment should use env vars instead of config file
		v.SetEnvPrefix("")
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	if v.IsSet("secrets") {
		sv := viper.New()
		sv.SetConfigFile(v.GetString("secrets"))
		if err := sv.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read secrets config: %w", err)
		}
		if err := v.MergeConfigMap(sv.AllSettings()); err != nil {
			return nil, fmt.Errorf("failed to merge secrets config: %w", err)
		}
	}

	return v, nil
}
