package dialect

import (
	"github.com/soderasen-au/go-common/util"
	"gorm.io/gorm"
	"strings"
)

type Dialect string

const (
	MYSQL     Dialect = "mysql"
	PG        Dialect = "postgresql"
	SQLSERVER Dialect = "sqlserver"
	SQLITE    Dialect = "sqlite"
)

type Opener interface {
	Open(dsn string) (*gorm.DB, error)
}

func Parse(str string) (*Dialect, Opener) {
	switch strings.ToLower(str) {
	case "mysql":
		return util.Ptr(MYSQL), &MysqlOpener{}
	case "postgresql", "postgres", "pg", "pgx":
		return util.Ptr(PG), &PostgresOpener{}
	case "sqlserver", "mssql":
		return util.Ptr(SQLSERVER), &SqlServerOpener{}
	case "sqlite", "sqlite3":
		return util.Ptr(SQLITE), &SQLiteOpener{}
	default:
		return nil, nil
	}
}
