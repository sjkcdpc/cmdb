package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/mds1455975151/cmdb/storage/model"
	"github.com/sirupsen/logrus"
	"fmt"
	"encoding/json"
)

var dbCmdb *gorm.DB

func CmdbDatabase() *gorm.DB {
	return dbCmdb
}

func InitCmdbDatabase() {

	//Migrate the schema 表结构定义
	dbCmdb.AutoMigrate(
		&model.Hosts{},
		&model.Users{},
		&model.Dnsinfo{},
		&model.Lbinfo{},
		&model.Logs{},
		&model.Businessinfo{},
	)

	// 初始化测试数据
	InsertTestData()
}

func QueryHost(id int64) *model.Hosts {

	var data model.Hosts

	if err := dbCmdb.Where("id = ?", id).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": id,
			"error":    err.Error(),
		}).Error("QueryHost data not found.")

		return nil
	}

	return &data
}

func QueryHostAll() interface{} {

	var n [] string
	datas := make([]*model.Hosts, 0)


	if err := dbCmdb.Find(&datas).Error; err !=nil{
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
		}).Error("QueryHostAll data not found.")
	}

	for _, data := range datas {
		json, _ := json.Marshal(data)
		n = append(n, string(json))
	}

	return n
}

// 获取表数据行数
func QueryHostsCount() int{
	var count int=0

	dbCmdb.Table("hosts").Count(&count)
	return count
}

func InsertTestData() *model.Hosts {

	//初始化数据
	host1 := model.Hosts{GlobalId:1,
							WanIp:"1.1.1.1",
							LanIp:"10.0.0.1",
							Conf:"4core+8G+50G+500G",
							HostName:"test1",
							Os:"CentOS 6",
							Manager:"dongsheng.ma",
							SshPort:22,
							Tags:"nginx",
							Remark:"测试数据"}
	host2 := model.Hosts{GlobalId:2,
						WanIp:"2.2.2.2",
						LanIp:"10.0.0.2",
						Conf:"4core+8G+50G+500G",
						HostName:"test2",
						Os:"CentOS 7",
						Manager:"dongsheng.ma",
						SshPort:22,
						Tags:"iis",
						Remark:"测试数据"}
	dbCmdb.Create(&host1)
	dbCmdb.Create(&host2)
	fmt.Println("insert test data")

	if err := dbCmdb.Create(&host1).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
		}).Error("Insert Test data not ok.")

		return nil
	}

	return nil
}

// 查询单个用户信息
func QueryUsers(id int64) *model.Users {

	var data model.Users

	if err := dbCmdb.Where("id = ?", id).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": id,
			"error":    err.Error(),
		}).Error("QueryUsers data not found.")

		return nil
	}

	return &data
}