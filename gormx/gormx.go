package gormx

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifeTime int
}

type DB struct {
	*gorm.DB
}

func New(dialect, dsn string, options ...Option) *DB {
	cfg := Config{
		MaxIdleConns:    10,
		MaxOpenConns:    10,
		ConnMaxLifeTime: 10,
	}

	for _, opt := range options {
		opt(&cfg)
	}

	ormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	var dialector gorm.Dialector
	switch dialect {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		if err := os.MkdirAll(filepath.Dir(dsn), os.ModePerm); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		dialector = sqlite.Open(dsn)
	default:
		log.Fatalf("Unsupported dialect: %s, supported dialects are - mysql, postgres, sqlite", dialect)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: ormLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeTime) * time.Second)
	return &DB{db}
}
