package db

import (
	"fmt"
	"github.com/soderasen-au/go-ipc/db/dialect"
	"gorm.io/gorm"
)

type Config struct {
	Dialect string `json:"dialect,omitempty" yaml:"dialect"`

	// Examples:
	//   - SQL Server: "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
	//   - MySQL: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	//   - PostgreSQL: user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	//   - SQLite: "file.db?_busy_timeout=5000"
	DSN string `json:"dsn,omitempty" yaml:"dsn"`
}

func NewDB(cfg Config) (*gorm.DB, error) {
	if d, opener := dialect.Parse(cfg.Dialect); d != nil {
		return opener.Open(cfg.DSN)
	}
	return nil, fmt.Errorf("DB dialect [%s] is not supported", cfg.Dialect)
}
