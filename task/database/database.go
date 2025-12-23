package database
import (
    "task/models"
    "log"
    "os"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
    // 从环境变量读取数据库配置
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        dsn = "user:password@tcp(localhost:3306)/contract_events?charset=utf8mb4&parseTime=True&loc=Local"
    }
	log.Println("Connecting to database with DSN:", dsn)

    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // 自动迁移表结构
    err = DB.AutoMigrate(&models.ContractEvent{})
    if err != nil {
        log.Fatal("Database migration failed:", err)
    }

    log.Println("Database connection established")
    return DB
}