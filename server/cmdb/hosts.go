package cmdb

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mds1455975151/cmdb/storage"
	//"github.com/sirupsen/logrus"
	"github.com/mds1455975151/cmdb/storage/model"
)

func init() {
	getHandlers["/hosts"] = hostsinfo
	handlers["/hosts"] = hostsinsert
}

func hostsinfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		Info := storage.QueryHostsAll()
		c.JSON(http.StatusOK, Info)
	} else {
		i, _ := strconv.ParseInt(id, 10, 64)
		Info := storage.QueryHost(i)
		c.JSON(http.StatusOK, Info)
	}
}

////////////////////////////////////////////////////////////////////////////////////// POST
func hostsinsert(c *gin.Context) {
	//fmt.Println(c.Request)
	wanip := c.Request.FormValue("wanip")
	lanip := c.Request.FormValue("lanip")
	conf := c.Request.FormValue("conf")
	hostname := c.Request.FormValue("hostname")
	os := c.Request.FormValue("os")
	contact := c.Request.FormValue("contact")
	manager := c.Request.FormValue("manager")
	tags := c.Request.FormValue("tags")
	remark := c.Request.FormValue("remark")

	//fmt.Println(wanip)
	//insert into hosts (hostname, conf, wan_ip, lan_ip, os) VALUES ("gate","2core*4G*20M*50G+300G","182.254.145.181","10.0.0.11","Tencent Linux Release 2.2 (Final)");
	hostsinsertinfo := &model.Hosts{
		Wanip:    wanip,
		Lanip:    lanip,
		Conf:     conf,
		Hostname: hostname,
		Os:       os,
		Contact:  contact,
		Manager:  manager,
		Tags:     tags,
		Remark:   remark,
	}
	//fmt.Println(hostsinfo)
	storage.Insert(storage.CmdbDatabase(), hostsinsertinfo)
	c.JSON(http.StatusOK, gin.H{"msg": "ok",})
}
