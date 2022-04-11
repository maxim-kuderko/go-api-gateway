package initializers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-api-gateway/entities"
	"os"
)

func RouterConfig(v *viper.Viper) *entities.Config {
	file, err := os.Open(v.GetString(`CONFIG_FILE`))
	if err != nil {
		logrus.Fatal(err)
	}
	var config *entities.Config
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		logrus.Fatal(err)
	}
	return config
}
