package gorm

import (
	"fmt"
	"github.com/gotechbook/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	Db *gorm.DB
}

func NewClient(dsn string, dev bool) *Client {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("gorm connect mysql err %v", err)
		return nil
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxIdleTime(20)
	if dev {
		db = db.Debug()
	}
	return &Client{Db: db}
}

func Dsn(username string, password string, host string, port int, db string, timeout string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, db, timeout)
}
