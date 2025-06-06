package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/xbmlz/pkgo/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

type Config struct {
	Driver          string `yaml:"driver" json:"driver" mapstructure:"driver" env:"DB_DRIVER"`
	Host            string `yaml:"host" json:"host" mapstructure:"host" env:"DB_HOST"`
	Port            string `yaml:"port" json:"port" mapstructure:"port" env:"DB_PORT"`
	Username        string `yaml:"username" json:"username" mapstructure:"username" env:"DB_USERNAME"`
	Password        string `yaml:"password" json:"password" mapstructure:"password" env:"DB_PASSWORD"`
	Database        string `yaml:"database" json:"database" mapstructure:"database" env:"DB_DATABASE"`
	Params          string `yaml:"params" json:"params" mapstructure:"params" env:"DB_PARAMS"`
	MaxIdleConns    int    `yaml:"max_idle_conns" json:"max_idle_conns" mapstructure:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns    int    `yaml:"max_open_conns" json:"max_open_conns" mapstructure:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	ConnMaxLifeTime int    `yaml:"conn_max_life_time" json:"conn_max_life_time" mapstructure:"conn_max_life_time" env:"DB_CONN_MAX_LIFE_TIME"`
}

func (c Config) DSN() string {
	switch c.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Database, c.Params)
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", c.Host, c.Port, c.Username, c.Password, c.Database, c.Params)
	case "sqlite":
		return c.Database
	default:
		return ""
	}
}

func New(cfg Config) *DB {
	ormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN())
	case "postgres":
		dialector = postgres.Open(cfg.DSN())
	case "sqlite":
		if err := os.MkdirAll(filepath.Dir(cfg.DSN()), os.ModePerm); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		dialector = sqlite.Open(cfg.DSN())
	default:
		log.Fatalf("Unsupported dialect: %s, supported dialects are - mysql, postgres, sqlite", cfg.DSN())
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
	sqlDB.SetMaxIdleConns(utils.OrElse(cfg.MaxIdleConns, 5))
	sqlDB.SetMaxOpenConns(utils.OrElse(cfg.MaxOpenConns, 10))
	sqlDB.SetConnMaxLifetime(time.Duration(utils.OrElse(cfg.ConnMaxLifeTime, 10)) * time.Second)
	return &DB{db}
}
