package mysql

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"book-management-system/configs"
)

var (
	mysqlDB *gorm.DB
	once    sync.Once
)

// Init returns mysqlDB connection instance
func Init() *gorm.DB {
	once.Do(func() {
		dsn := getMySQLConnString()
		fmt.Println(dsn)

		var err error
		mysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to mysql database: %s", err)
		}
	})

	return mysqlDB
}

// getMySQLConnString return connection string from config
func getMySQLConnString() string {
	config := configs.GetConfig()
	c := config.Mysql

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.User, c.Pass, c.Host, c.Port, c.Name)
}
