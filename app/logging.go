package app

import "github.com/Sirupsen/logrus"

var LOG = logrus.New()

func InitializeLogging(config Config) {
	if config.Debug {
		LOG.Level = logrus.DebugLevel
	}
}
