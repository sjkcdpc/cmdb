package storage

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/mds1455975151/cmdb/settings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Initialize() {

	dbCmdb = InitMySql(settings.Get("cmdb"))
	InitCmdbDatabase()
}

func InitMySql(setting *viper.Viper) *gorm.DB {

	var db *gorm.DB

	if setting == nil {
		logrus.Fatalf("GORM InitMySql failed. No setting specified.")
	}

	var host = setting.GetString("mysql.host")
	var port = setting.GetInt64("mysql.port")
	var user = setting.GetString("mysql.user")
	var password = setting.GetString("mysql.password")
	var database = setting.GetString("mysql.database")

	var createDatabaseIfNotExist = setting.GetBool("mysql.create_database_if_not_exist")

	var connString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	logrus.Infof("Connection to MySql: %s", connString)

	//open a db connection
	var err error
	db, err = gorm.Open("mysql", connString)
	if err != nil {

		Error1049 := strings.Contains(err.Error(), "1049")

		if !Error1049 || !createDatabaseIfNotExist {
			logrus.WithFields(logrus.Fields{
				"connString": connString,
				"error":      err.Error(),
			}).Fatalf("Connection mysql failed. If database is not exist, create database with the command:"+
				"'CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'", database)
		}

		createDatabase(setting)

		db, err = gorm.Open("mysql", connString)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"connString": connString,
			"error":      err.Error(),
		}).Fatalf("Connection mysql failed. If database is not exist, create database with the command:"+
			"'CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'", database)
	}

	return db
}

func createDatabase(setting *viper.Viper) {

	var host = setting.GetString("mysql.host")
	var port = setting.GetInt64("mysql.port")
	var user = setting.GetString("mysql.user")
	var password = setting.GetString("mysql.password")
	var database = setting.GetString("mysql.database")

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8", user, password, host, port)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"dataSourceName": dataSourceName,
			"error":          err,
		}).Error("cmdb mysql connection failed.")
	}
	defer db.Close()

	createDatabaseSql := fmt.Sprintf(
		`CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`,
		database)

	// create database
	rows, err := db.Query(createDatabaseSql)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dataSourceName":    dataSourceName,
			"createDatabaseSql": createDatabaseSql,
			"rows":              rows,
			"error":             err,
		}).Error("cmdb mysql create database failed.")
	}
	rows.Close()

	logrus.WithFields(logrus.Fields{
		"dataSourceName":    dataSourceName,
		"createDatabaseSql": createDatabaseSql,
		"rows":              rows,
	}).Warn("cmdb create database success.")
}

func Insert(db *gorm.DB, model interface{}) {

	if db == nil {

		logrus.WithFields(logrus.Fields{
			"db":    db,
			"model": model,
		}).Error("database unavailable")

		return
	}

	db.Create(model)
}

func Save(db *gorm.DB, model interface{}) {

	if db == nil {

		logrus.WithFields(logrus.Fields{
			"db":    db,
			"model": model,
		}).Error("database unavailable")

		return
	}

	db.Save(model)
}
