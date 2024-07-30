package dialect

import (
	"strings"

	"gorm.io/gorm"

	"github.com/soderasen-au/go-common/util"
)

type DBType string

const (
	MYSQL     DBType = "mysql"
	PG        DBType = "postgresql"
	SQLSERVER DBType = "sqlserver"
	SQLITE    DBType = "sqlite"
)

type Opener interface {
	Open(dsn string) (*gorm.DB, error)
}

func Parse(str string) (*DBType, Opener) {
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
