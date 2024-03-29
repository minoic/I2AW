package database

import (
	"github.com/MinoIC/I2AW/configure"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func connect() {
	conf := configure.GetConf()
	dialect := conf.String("database")
	switch dialect {
	case "SQLITE":
		DB, err := gorm.Open("sqlite3", "sqlite3.db")
		if err != nil {
			db = nil
			beego.Error(err.Error())
		}
		db = DB
		return
	case "MYSQL":
		DSN := conf.String("MYSQLUsername") + ":" +
			conf.String("MYSQLUserPassword") + "@(" +
			conf.String("MYSQLHost") + ")/" +
			conf.String("MYSQLDatabaseName") +
			"?charset=utf8&parseTime=True&loc=Local"
		DB, err := gorm.Open("mysql", DSN)
		if err != nil {
			db = nil
			beego.Error(err.Error())
		}
		db = DB
		return
	}
	db = nil
	panic("CONF ERR: WRONG SQL DIALECT!!! " + dialect)
}

func GetDatabase() *gorm.DB {
	for db == nil {
		beego.Warn("trying to connect to database!")
		connect()
	}

	for err := db.DB().Ping(); err != nil; err = db.DB().Ping() {
		beego.Warn("trying to connect to database!")
		connect()
	}
	return db
}
