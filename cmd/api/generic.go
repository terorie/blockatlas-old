package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas"
	"net/http"
)

func makeTxRoute(router gin.IRouter, api blockatlas.TxAPI) {
	router.GET("/:address/txs", func(c *gin.Context) {
		address := c.Param("address")
		if address == "" {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		page, err := api.GetTxsByAddress(address)
		if !handleError(c, err) {
			return
		}

		page.Sort()
		c.JSON(http.StatusOK, &page)
	})
}

func makeTokenTxRoute(router gin.IRouter, api blockatlas.TokenTxAPI) {
	router.GET("/:address/token/:token/txs", func(c *gin.Context) {
		address := c.Param("address")
		token := c.Param("token")
		if address == "" || token == "" {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		page, err := api.GetTokenTxsByAddress(address, token)
		if !handleError(c, err) {
			return
		}

		page.Sort()
		c.JSON(http.StatusOK, &page)
	})
}

func handleError(c *gin.Context, err error) bool {
	switch {
	case err == blockatlas.ErrInvalidAddr:
		c.String(http.StatusBadRequest, "Invalid address")
		return false
	case err == blockatlas.ErrNotFound:
		c.String(http.StatusNotFound, "No such address")
		return false
	case err == blockatlas.ErrSourceConn:
		c.String(http.StatusServiceUnavailable, "Lost connection to blockchain")
		return false
	case err != nil:
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return false
	}
	return true
}
