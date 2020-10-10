package mysql

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func setupTestSuite() (*gorm.DB, sqlmock.Sqlmock, error) {
	dbMock, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		return nil, nil, err
	}

	gormDBMock, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      dbMock,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DisableAutomaticPing: true,
		Logger:               logger.Discard,
	})
	if err != nil {
		return nil, nil, err
	}

	return gormDBMock, mock, nil
}

func closeDB(gormDBMock *gorm.DB) {
	dbMock, _ := gormDBMock.DB()
	dbMock.Close()
}

var errDatabase = errors.New("error")
