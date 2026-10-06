package db

import (
	"fmt"
	"os"

	"github.com/SahilMahale/notes-backend/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConnection struct {
	Db *gorm.DB
}

func NewDbConnection() (DbConnection, error) {
	host := utils.GetEnvDefault("DB_HOST", "localhost")
	user := utils.GetEnvDefault("DB_USER", "noto")
	password := utils.GetEnvDefault("DB_PASSWORD", "noto_dev_password")
	dsn := makeDSN(host, user, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return DbConnection{}, fmt.Errorf("Error establishing DB connection: %v", err)
	}
	db.AutoMigrate(&Note{})
	db.AutoMigrate(&User{})
	return DbConnection{Db: db}, nil
}

func makeDSN(host, user, password string) string {
	sslMode := "disable"
	dbEnv := os.Getenv("DB_ENV")
	dsnTemplate := "host=%s user=%s password=%s dbname=noto port=5432 sslmode=%s TimeZone=%s"
	if dbEnv != "dev" {
		sslMode = "require"
	}
	tz := utils.SystemTimeZoneName()
	return fmt.Sprintf(dsnTemplate, host, user, password, sslMode, tz)
}
