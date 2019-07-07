package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadConfig(confPath string) {
	loadDefaults()

	// Load config from environment
	viper.SetEnvPrefix("atlas")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// Load config file
	if confPath == "" {
		err := viper.ReadInConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Info("Running without config file")
		} else if err != nil {
			logrus.Fatal(err)
		} else {
			logrus.Info("Using config.yml")
		}
	} else {
		viper.SetConfigFile(confPath)
		err := viper.ReadInConfig()
		if err != nil {
			logrus.WithError(err).Error("Failed to read config")
		} else {
			logrus.WithField("config_file", confPath).Info("Using config file")
		}
	}
}

func loadDefaults() {
	viper.SetDefault("gin.mode", gin.ReleaseMode)
	viper.SetDefault("gin.reverse_proxy", false)
	viper.SetDefault("observer.redis", "redis://localhost:6379")
	viper.SetDefault("observer.min_poll", 250 * time.Millisecond)
	viper.SetDefault("observer.backlog", 3 * time.Hour)
	viper.SetDefault("observer.backlog_max_blocks", 200)
	viper.SetDefault("observer.stream_conns", 16)

	// All platforms with public RPC endpoints
	viper.SetDefault("binance.api", "https://explorer.binance.org/api/v1")
	viper.SetDefault("ripple.api", "https://data.ripple.com/v2")
	viper.SetDefault("stellar.api", "https://horizon.stellar.org")
	viper.SetDefault("kin.api", "https://horizon.kinfederation.com/")
	viper.SetDefault("tezos.api", "https://api1.tzscan.io/v3")
	viper.SetDefault("aion.api", "https://mainnet-api.aion.network/aion/dashboard")
	viper.SetDefault("icon.api", "https://tracker.icon.foundation/v3")
	viper.SetDefault("vechain.api", "https://explore.veforge.com/api")
	viper.SetDefault("theta.api", "https://explorer.thetatoken.org:9000/api")
	viper.SetDefault("cosmos.api", "https://stargate.cosmos.network")
	viper.SetDefault("semux.api", "https://sempy.online/api")
	viper.SetDefault("ontology.api", "https://explorer.ont.io/api/v1/explorer")
	viper.SetDefault("iotex.api", "https://pharos.iotex.io/v1")
	viper.SetDefault("waves.api", "https://nodes.wavesnodes.com")
}
