package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/mds1455975151/cmdb/storage/model"
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

