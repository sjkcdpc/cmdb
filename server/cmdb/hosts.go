package cmdb

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mds1455975151/cmdb/storage"
)

func init() {
	getHandlers["/id"] = hostsone
	getHandlers["/count"] = hostscount
	getHandlers["/all"] = hostsall
}

func hostsone(c *gin.Context) {
	info := storage.QueryHost(1)
	c.JSON(http.StatusOK, info)
}

func hostscount(c *gin.Context) {
	info := storage.QueryHostsCount()
	c.JSON(http.StatusOK, info)
}

func hostsall(c *gin.Context) {
	info := storage.QueryHostAll()
	c.JSON(http.StatusOK, info)
}