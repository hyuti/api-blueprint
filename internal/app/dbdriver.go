package app

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsnTemplate = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai"

func dsnBuilder(host, port, user, pwd, dbName string) string {
	return fmt.Sprintf(dsnTemplate, host, port, user, pwd, dbName)
}

func WithDBDriver() (*gorm.DB, error) {
	if app.cfg == nil {
		return nil, ErrCfgEmpty
	}
	db, err := gorm.Open(postgres.Open(dsnBuilder(
		app.cfg.Postgres.Host,
		app.cfg.Postgres.Port,
		app.cfg.Postgres.User,
		app.cfg.Postgres.Password,
		app.cfg.Postgres.DatabaseName)))
	if err != nil {
		return nil, fmt.Errorf("cannot to init db driver: %w", err)
	}
	return db, nil
}

func DBDriver() *gorm.DB {
	mutex.Lock()
	defer mutex.Unlock()
	return app.dbDriver
}
