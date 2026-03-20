package database

import (
	"fmt"
	"log"
	"my-blog-project/config"
	"my-blog-project/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	// 从配置文件中读取数据库的信息
	con, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置文件失败")
	}
	dbHost := con.Database.Host
	dbPort := con.Database.Port
	dbUsername := con.Database.Username
	dbPassword := con.Database.Password
	dbName := con.Database.DBName

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to MySQL database:", err)
	}

	err = DB.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("MySQL database connected and migrated successfully")
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}
