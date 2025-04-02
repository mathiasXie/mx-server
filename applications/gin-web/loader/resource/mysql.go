package resource

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/config"
	log "github.com/mathiasXie/gin-web/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxIdl  = 20
	maxOpen = 50
)

func InitDB(mysqlConfig *config.MysqlConfig) (*gorm.DB, *gorm.DB) {
	if mysqlConfig.Host == "" {
		return nil, nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DBName, mysqlConfig.Charset)

	gormLogger := log.NewGormLogger()
	if mysqlConfig.LogLevel > 0 {
		gormLogger.LogLevel = logger.LogLevel(mysqlConfig.LogLevel)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.CtxError(&gin.Context{}, "[InitDB]Error", err.Error())
		panic("failed to connect database")
	}
	fmt.Println("connect mysql success")
	return db, db
}
