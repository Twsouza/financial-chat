package data

import (
	"log"
	"os"
	"server/core"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Db            *gorm.DB
	Logger        *log.Logger
	Dsn           string
	DsnTest       string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "test"
	dbInstance.DsnTest = "file::memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debug = true

	connection, err := dbInstance.Connect()
	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (db *Database) Connect() (*gorm.DB, error) {
	var err error

	config := &gorm.Config{}
	if db.Debug {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				Colorful:                  false,
			},
		)

		config.Logger = newLogger
	}

	if db.Env != "test" {
		db.Db, err = gorm.Open(postgres.Open(db.Dsn))
	} else {
		db.Db, err = gorm.Open(sqlite.Open(db.DsnTest), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	if db.AutoMigrateDb {
		db.Db.AutoMigrate(core.User{})
	}

	return db.Db, nil
}
