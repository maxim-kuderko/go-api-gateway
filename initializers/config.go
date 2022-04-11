package initializers

import "github.com/spf13/viper"

func NewConfig() *viper.Viper {
	v := viper.New()
	v.SetDefault(`CONFIG_FILE`, `./config.json`)
	v.AutomaticEnv()
	return v
}
