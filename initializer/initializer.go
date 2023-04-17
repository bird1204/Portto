package initializer

import (
	"fmt"

	"portto/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	DB  *gorm.DB
	GIN *gin.Engine
}

func (s *Server) InitializeDB(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	// create database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE DATABASE IF NOT EXISTS " + DbName)

	// create tables
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	s.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	err = s.DB.Debug().Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.Block{},
		&model.Transaction{},
	)
	if err != nil {
		panic(err.Error())
	}
}

func (s *Server) InitializeGin() {
	s.GIN = gin.Default()
}
