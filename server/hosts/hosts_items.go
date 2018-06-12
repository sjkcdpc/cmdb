package hosts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mds1455975151/cmdb/storage"
)

func init() {

	handlerFuncList["/all"] = hostsall
}



func hostsall(c *gin.Context) {
	info := storage.QueryHost(1)
	c.JSON(http.StatusOK, info)
}