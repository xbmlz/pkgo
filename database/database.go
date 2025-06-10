package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	oracle "github.com/godoes/gorm-oracle"
	"github.com/xbmlz/pkgo/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DB struct {
	*gorm.DB
}

type Config struct {
	Driver          string `yaml:"driver" json:"driver" mapstructure:"driver" env:"DB_DRIVER"`
	Host            string `yaml:"host" json:"host" mapstructure:"host" env:"DB_HOST"`
	Port            int    `yaml:"port" json:"port" mapstructure:"port" env:"DB_PORT"`
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
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Database, c.Params)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s", c.Host, c.Port, c.Username, c.Password, c.Database, c.Params)
	case "sqlite":
		return c.Database
	case "mssql":
		// "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", c.Username, c.Password, c.Host, c.Port, c.Database)
	case "oracle":
		// oracle://user:password@127.0.0.1:1521/service
		options := map[string]string{
			"CONNECTION TIMEOUT": "90",
			"LANGUAGE":           "SIMPLIFIED CHINESE",
			"TERRITORY":          "CHINA",
			"SSL":                "false",
		}
		return oracle.BuildUrl(c.Host, c.Port, c.Database, c.Username, c.Password, options)
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
	case "mssql":
		dialector = sqlserver.Open(cfg.DSN())
	case "oracle":

		dialector = oracle.New(oracle.Config{
			DSN:                     cfg.DSN(),
			IgnoreCase:              false, // query conditions are not case-sensitive
			NamingCaseSensitive:     true,  // whether naming is case-sensitive
			VarcharSizeIsCharLength: true,  // whether VARCHAR type size is character length, defaulting to byte length

			// RowNumberAliasForOracle11 is the alias for ROW_NUMBER() in Oracle 11g, defaulting to ROW_NUM
			RowNumberAliasForOracle11: "ROW_NUM",
		})
	default:
		log.Fatalf("Unsupported dialect: %s, supported dialects are - mysql, postgres, sqlite, mssql, oracle", cfg.DSN())
	}

	dbConfig := &gorm.Config{
		Logger: ormLogger,
	}

	if cfg.Driver == "oracle" {
		// 是否禁用默认在事务中执行单次创建、更新、删除操作
		dbConfig.SkipDefaultTransaction = true
		// 是否禁止在自动迁移或创建表时自动创建外键约束
		dbConfig.DisableForeignKeyConstraintWhenMigrating = true
		// 自定义命名策略
		dbConfig.NamingStrategy = schema.NamingStrategy{
			NoLowerCase:         true, // 是否不自动转换小写表名
			IdentifierMaxLength: 30,   // Oracle: 30, PostgreSQL:63, MySQL: 64, SQL Server、SQLite、DM: 128
		}
		// 创建并缓存预编译语句，启用后可能会报 ORA-01002 错误
		dbConfig.PrepareStmt = false
		// 插入数据默认批处理大小
		dbConfig.CreateBatchSize = 50
	}

	db, err := gorm.Open(dialector, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	if cfg.Driver == "oracle" {
		_, _ = oracle.AddSessionParams(sqlDB, map[string]string{
			"TIME_ZONE":               "+08:00",                       // ALTER SESSION SET TIME_ZONE = '+08:00';
			"NLS_DATE_FORMAT":         "YYYY-MM-DD",                   // ALTER SESSION SET NLS_DATE_FORMAT = 'YYYY-MM-DD';
			"NLS_TIME_FORMAT":         "HH24:MI:SSXFF",                // ALTER SESSION SET NLS_TIME_FORMAT = 'HH24:MI:SS.FF3';
			"NLS_TIMESTAMP_FORMAT":    "YYYY-MM-DD HH24:MI:SSXFF",     // ALTER SESSION SET NLS_TIMESTAMP_FORMAT = 'YYYY-MM-DD HH24:MI:SS.FF3';
			"NLS_TIME_TZ_FORMAT":      "HH24:MI:SS.FF TZR",            // ALTER SESSION SET NLS_TIME_TZ_FORMAT = 'HH24:MI:SS.FF3 TZR';
			"NLS_TIMESTAMP_TZ_FORMAT": "YYYY-MM-DD HH24:MI:SSXFF TZR", // ALTER SESSION SET NLS_TIMESTAMP_TZ_FORMAT = 'YYYY-MM-DD HH24:MI:SS.FF3 TZR';
		})
	}
	sqlDB.SetMaxIdleConns(utils.OrElse(cfg.MaxIdleConns, 5))
	sqlDB.SetMaxOpenConns(utils.OrElse(cfg.MaxOpenConns, 10))
	sqlDB.SetConnMaxLifetime(time.Duration(utils.OrElse(cfg.ConnMaxLifeTime, 10)) * time.Second)
	return &DB{db}
}
