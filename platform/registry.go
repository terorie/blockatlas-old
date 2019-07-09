package platform

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas"
	"github.com/trustwallet/blockatlas/platform/nimiq"
)

var platformList = []blockatlas.Platform{
	&nimiq.Platform{},
}

// Platforms contains all registered platforms by handle
var Platforms map[string]blockatlas.Platform

func Init() {
	Platforms = make(map[string]blockatlas.Platform)

	for _, platform := range platformList {
		handle := platform.Coin().Handle
		apiKey := fmt.Sprintf("%s.api", handle)

		if viper.GetString(apiKey) == "" {
			continue
		}

		log := logrus.WithFields(logrus.Fields{
			"platform": handle,
			"coin": platform.Coin(),
		})

		if _, exists := Platforms[handle]; exists {
			log.Fatal("Duplicate handle")
		}

		err := platform.Init()
		if err != nil {
			log.WithError(err).Fatal("Failed to initialize API")
		}

		Platforms[handle] = platform
	}
}
