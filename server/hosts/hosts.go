package hosts

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	//"github.com/sirupsen/logrus"
)
type hosts struct {
	id     string `form:"id" json:"id" binding:"required"`
}

func Start(group *gin.RouterGroup) {

	group.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})


	//logrus.WithFields(logrus.Fields{
	//	"n":     n,
	//	"error": err,
	//}).Error("exec prefix")

	group.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	group.POST("/post", func(c *gin.Context) {

		id := c.PostForm("id")
		page := c.PostForm("page")
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
		messages := "id:"
		messages += id
		messages += ",page:"
		messages += page
		messages += ",name:"
		messages += name
		messages += ",message:"
		messages += message
		c.String(http.StatusOK, messages)
	})
}