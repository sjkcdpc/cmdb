package server

import (
	"github.com/sirupsen/logrus"
	//"github.com/mds1455975151/cmdb/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/mds1455975151/cmdb/server/cmdb"
	"github.com/mds1455975151/cmdb/settings"
	"github.com/mds1455975151/cmdb/storage"
	"os"
)

var logo = `
		__                               __  __
		/ /___ _   _____     ____ ___  __/ /_/ /_
		/ / __ \ | / / _ \   / __ '/ / / / __/ __ \
		/ / /_/ / |/ /  __/  / /_/ / /_/ / /_/ / / /
		/_/\____/|___/\___/   \__,_/\__,_/\__/_/ /_/
		BasePath: /api/v1
		Version: v1.0.0 -- develop(00a2ad5)
		Contact: mds1455975151<1455975151@qq.com>
`

func Run() {
	logrus.Info(logo)

	//log.InitLogger()
	storage.Initialize()

	var setting = settings.Get("cmdb")
	var listen = setting.GetString("gin.listen")
	var mode = setting.GetString("gin.mode")

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add a ginrus middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	//router.Use(utils.Ginrus(logrus.StandardLogger(), time.RFC3339, true))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// v1，v2 分组后期可以扩展
	var basePath = "/api/v1"

	v1 := router.Group(basePath)
	{
		cmdb.Start(v1.Group("/cmdb"))
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//noinspection SpellCheckingInspection
	logrus.WithFields(logrus.Fields{
		"pid":     os.Getpid(),
		"listen":  listen,
		"version": "v1.0.0 -- develop(00a2ad5)",
		"path":    basePath,
	}).Info("Start cmdb service.")

	router.Run(listen)
}
