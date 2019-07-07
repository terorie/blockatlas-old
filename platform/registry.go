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

// BlockAPIs contains platforms with block services
var BlockAPIs map[string]blockatlas.BlockAPI

// CustomAPIs contains platforms with custom HTTP services
var CustomAPIs map[string]blockatlas.CustomAPI

func Init() {
	Platforms  = make(map[string]blockatlas.Platform)
	BlockAPIs  = make(map[string]blockatlas.BlockAPI)
	CustomAPIs = make(map[string]blockatlas.CustomAPI)

	for _, platform := range platformList {
		handle := platform.Coin().Handle
		apiKey := fmt.Sprintf("%s.api", handle)

		if !viper.IsSet(apiKey) {
			continue
		}
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

		if blockAPI, ok := platform.(blockatlas.BlockAPI); ok {
			BlockAPIs[handle] = blockAPI
		}
		if customAPI, ok := platform.(blockatlas.CustomAPI); ok {
			CustomAPIs[handle] = customAPI
		}
	}
}
