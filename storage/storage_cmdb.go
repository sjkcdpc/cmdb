package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/mds1455975151/cmdb/storage/model"
	"github.com/sirupsen/logrus"
)

var dbCmdb *gorm.DB

func CmdbDatabase() *gorm.DB {
	return dbCmdb
}

func InitAuthDatabase() {

	//Migrate the schema 表结构定义
	dbCmdb.AutoMigrate(
		&model.Hosts{},
	)
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