package dbconfig

import (
	"fmt"

	"github.com/MyriadFlow/storefront-gateway/config/envconfig"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Return singleton instance of db, initiates it before if it is not initiated already
func GetDb() *gorm.DB {
	if db != nil {
		return db
	}
	var (
		host     = envconfig.EnvVars.DB_HOST
		username = envconfig.EnvVars.DB_USERNAME
		password = envconfig.EnvVars.DB_PASSWORD
		dbname   = envconfig.EnvVars.DB_NAME
		port     = envconfig.EnvVars.DB_PORT
	)

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%d",
		host, username, password, dbname, port)

	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dns,
	}))
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal("failed to ping database", err)
	}
	if err = sqlDb.Ping(); err != nil {
		log.Fatal("failed to ping database", err)
	}
	return db
}
