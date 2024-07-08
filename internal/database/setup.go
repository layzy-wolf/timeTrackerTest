package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/layzy-wolf/timeTrackerTest/internal/env"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

import _ "github.com/golang-migrate/migrate/source/file"

func Setup(conf *env.Config) *sqlx.DB {

	log.Debugln("connect to database")
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.Database))
	if err != nil {
		log.Fatalln(err)
	}

	log.Debugln("migrations up")
	if err := driverMigrate(db.DB); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Errorln(err)
	}

	return db
}

func driverMigrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Debugln("error while creating driver instance")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/database/migrations", "postgres", driver)
	if err != nil {
		log.Debugln("error while creating migration instance")
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
