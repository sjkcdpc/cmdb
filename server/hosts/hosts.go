package hosts

import (
	//"net/http"
	//"fmt"

	"github.com/gin-gonic/gin"
	//"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus"
)
type hosts struct {
	id     string `form:"id" json:"id" binding:"required"`
}

var handlerFuncList = make(map[string]gin.HandlerFunc)

func Start(group *gin.RouterGroup) {

	if group == nil {
		logrus.Error("admin start failed.")
		return
	}

	for key, value := range handlerFuncList {
		group.GET(key, value)
	}
}