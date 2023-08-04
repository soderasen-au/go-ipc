package dialect

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SqlServerOpener struct{}

func (SqlServerOpener) Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
}
