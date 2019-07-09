package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/trustwallet/blockatlas"
	"github.com/trustwallet/blockatlas/platform"
)

var routers = make(map[string]gin.IRouter)

func loadPlatforms(root gin.IRouter) {
	v2 := root.Group("/v2")
	for handle, p := range platform.Platforms {
		loadPlatform(v2, handle, p)
	}

	logrus.WithField("routes", len(routers)).
		Info("Routes set up")

	v2.GET("/", getEnabledEndpoints)
}

func loadPlatform(g gin.IRouter, handle string, p blockatlas.Platform) {
	if customAPI, ok := p.(blockatlas.CustomAPI); ok {
		customAPI.RegisterRoutes(getRouter(g, handle))
	}
	if txAPI, ok := p.(blockatlas.TxAPI); ok {
		makeTxRoute(getRouter(g, handle), txAPI)
	}
	if tokenTxAPI, ok := p.(blockatlas.TokenTxAPI); ok {
		makeTokenTxRoute(getRouter(g, handle), tokenTxAPI)
	}
}

// getRouter lazy loads routers
func getRouter(router gin.IRouter, handle string) gin.IRouter {
	if group, ok := routers[handle]; ok {
		return group
	} else {
		path := fmt.Sprintf("/%s", handle)
		logrus.Debugf("Registering %s", path)
		group := router.Group(path)
		routers[handle] = group
		return group
	}
}
