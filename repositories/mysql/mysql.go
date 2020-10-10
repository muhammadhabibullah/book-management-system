// Package mysql contains all MySQL repositories
package mysql

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"book-management-system/configs"
	"book-management-system/entities/models"
)

var (
	mysqlDB *gorm.DB
	once    sync.Once
)

// Init returns mysqlDB connection instance
func Init() *gorm.DB {
	once.Do(func() {
		cfg := configs.GetConfig()
		dsn := getMySQLConnString(cfg.Mysql)

		var err error
		mysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to mysql database: %s", err)
		}

		if !cfg.Production {
			if err = mysqlDB.Set("gorm:table_options", "ENGINE=InnoDB").
				AutoMigrate(
					&models.Book{},
					&models.Member{},
				); err != nil {
				log.Fatalf("failed to migrate new model to mysql database: %s", err)
			}
		}

		configMySQLConn(cfg.Mysql)
	})

	return mysqlDB
}

// getMySQLConnString return connection string from config
func getMySQLConnString(cfg configs.MySQLConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
}

// configMySQLConn configure MySQLConnection settings
func configMySQLConn(cfg configs.MySQLConfig) {
	db, _ := mysqlDB.DB()
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
}
