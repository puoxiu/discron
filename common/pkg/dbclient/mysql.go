package dbclient

import (
	"github.com/puoxiu/discron/common/models"
	"github.com/puoxiu/discron/common/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _defaultDB *gorm.DB


func Init(m models.Mysql) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN: m.Dsn(),
		DefaultStringSize: 256,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{}); err!= nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		_defaultDB=db
		return db
	}
}


func GetMySQLDB() *gorm.DB {
	if _defaultDB==nil{
		logger.Errorf("mysql database is not initialized")
		return nil
	}
	return _defaultDB
}