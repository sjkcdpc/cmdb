package cmdb

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)
type hosts struct {
	id     string `form:"id" json:"id" binding:"required"`
}

var handlers = make(map[string]gin.HandlerFunc)
var getHandlers = make(map[string]gin.HandlerFunc)

func Start(group *gin.RouterGroup) {

	if group == nil {
		logrus.Error("admin start failed.")
		return
	}
	for key, value := range handlers {
		group.POST(key, value)
	}

	for key, value := range getHandlers {
		group.GET(key, value)
	}
}