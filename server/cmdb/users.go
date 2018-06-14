package cmdb

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mds1455975151/cmdb/storage"
)

func init() {
	getHandlers["/users"] = users
}

func users(c *gin.Context) {
	id :=c.Query("id")
	if id == ""{
		Info := storage.QueryUsersAll()
		c.JSON(http.StatusOK, Info)
	} else {
		i, _ := strconv.ParseInt(id, 10, 64)
		Info := storage.QueryUsers(i)
		c.JSON(http.StatusOK, Info)
	}
}
