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
		dsn := getMySQLConnString()

		var err error
		mysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to mysql database: %s", err)
		}

		if err = mysqlDB.Set("gorm:table_options", "ENGINE=InnoDB").
			AutoMigrate(
				&models.Book{},
			); err != nil {
			log.Fatalf("failed to migrate new model to mysql database: %s", err)
		}

		db, _ := mysqlDB.DB()
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(5 * time.Minute)
	})

	return mysqlDB
}

// getMySQLConnString return connection string from config
func getMySQLConnString() string {
	config := configs.GetConfig()
	c := config.Mysql

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Pass, c.Host, c.Port, c.Name)
}
