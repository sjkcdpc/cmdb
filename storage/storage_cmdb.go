package storage

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/mds1455975151/cmdb/storage/model"
	"github.com/sirupsen/logrus"
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
	//InsertHostTestData()
	//InsertUserTestData()
}

// 获取单个主机信息
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

// 获取所有主机信息
func QueryHostsAll() interface{} {

	var n []string
	datas := make([]*model.Hosts, 0)

	if err := dbCmdb.Find(&datas).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryHostAll data not found.")
	}

	for _, data := range datas {
		json, _ := json.Marshal(data)
		n = append(n, string(json))
	}

	return n
}

// 获取表数据行数
func QueryHostsCount() int {
	var count int = 0

	dbCmdb.Table("hosts").Count(&count)
	return count
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

// 获取所有用户信息
func QueryUsersAll() interface{} {

	var n []string
	datas := make([]*model.Users, 0)

	if err := dbCmdb.Find(&datas).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryUsersAll data not found.")
	}

	for _, data := range datas {
		json, _ := json.Marshal(data)
		n = append(n, string(json))
	}

	return n
}

// 获取表数据行数
func QueryUsersCount() int {
	var count int = 0

	dbCmdb.Table("users").Count(&count)
	return count
}