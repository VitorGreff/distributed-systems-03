package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPsqlConn() (*gorm.DB, error) {
	psqlInfo := "host=localhost user=postgres password=1234 dbname=sd2 port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
