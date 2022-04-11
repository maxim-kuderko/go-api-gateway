package initializers

import (
	"github.com/sirupsen/logrus"
	"os"
)

const defaultLogLevel = logrus.WarnLevel
const logLevelEnvVar = `LOG_LEVEL`

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	lvl := defaultLogLevel
	if v := os.Getenv(logLevelEnvVar); v != `` {
		if tmp, err := logrus.ParseLevel(v); err == nil {
			lvl = tmp
		}
	}
	logrus.SetLevel(lvl)
}
