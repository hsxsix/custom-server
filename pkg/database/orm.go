/**
 * @File: orm.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 10:03 AM
 */

package database

import (
	"custom_server/pkg/config"
	"custom_server/pkg/log"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"time"
)

func InitDB(cfg *config.DataBase) (*gorm.DB, error) {
	sqlLogger := log.NewGormLogger(log.GormLogConfig{
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  log.ParseLevel(cfg.LogLevel),
	})

	var (
		openDB *gorm.DB
		err    error
	)

	switch strings.ToUpper(cfg.Type) {
	case "MYSQL", "MARIADB":
		openDB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
			Logger: sqlLogger,
		})
	case "POSTGRES":
		openDB, err = gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
			Logger: sqlLogger,
		})
	case "SQLITE":
		openDB, err = gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{
			Logger: sqlLogger,
		})
	default:
		return nil, errors.New("does not support databases")
	}

	sqlDB, err := openDB.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxConn)
	sqlDB.SetMaxIdleConns(cfg.MaxConn / 2)
	return openDB, nil
}
