package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySqlConn() (*gorm.DB, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(err)
		log.Fatal("Couldnt get environment variables")
	}
	db_user := os.Getenv("MYSQL_USER")
	db_pass := os.Getenv("MYSQL_PASSWORD")
	db_host := os.Getenv("MYSQL_HOST")
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/sd?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_pass, db_host)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return nil, errors.New("erro ao conectar com banco de produtos")
	}

	return db, nil
}
