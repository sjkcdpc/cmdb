package storage

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/mds1455975151/cmdb/storage/model"
)

func InsertHostTestData() *model.Hosts {

	//初始化数据
	host1 := model.Hosts{Wanip: "1.1.1.1",
		Lanip: "10.0.0.1",
		Conf: "4core+8G+50G+500G",
		Hostname: "test1",
		Os: "CentOS 6",
		Manager: "dongsheng.ma",
		Sshport: 22,
		Tags: "nginx",
		Remark: "测试数据"}
	host2 := model.Hosts{Wanip: "2.2.2.2",
		Lanip: "10.0.0.2",
		Conf: "4core+8G+50G+500G",
		Hostname: "test2",
		Os: "CentOS 7",
		Manager: "dongsheng.ma",
		Sshport: 22,
		Tags: "iis",
		Remark: "测试数据"}
	dbCmdb.Create(&host1)
	dbCmdb.Create(&host2)
	fmt.Println("insert test data")

	if err := dbCmdb.Create(&host1).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Insert Test data not ok.")

		return nil
	}

	return nil
}

func InsertUserTestData() *model.Users {

	//初始化数据
	user1 := model.Users{
		Username: "madongsheng",
		Password: "123456",
		Birthday: "1991.05.15",
		Email:    "dongsheng.ma@lemongrassmedia.cn",
		Phone:    13693645328,
		Department:  "系统管理部",
		Enable:  true,
		Tags:     "试用期",
		Remark:   "测试数据"}
	user2 := model.Users{
		Username: "zhangsan",
		Password: "123456",
		Birthday: "1991.05.15",
		Email:    "zhangsan@lemongrassmedia.cn",
		Phone:    11111111111,
		Department:  "系统管理部",
		Enable:  false,
		Tags:     "转正",
		Remark:   "测试数据"}
	dbCmdb.Create(&user1)
	dbCmdb.Create(&user2)
	fmt.Println("users insert test data")

	if err := dbCmdb.Create(&user1).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("users insert test data not ok.")

		return nil
	}

	return nil
}
