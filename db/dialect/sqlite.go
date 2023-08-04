package dialect

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteOpener struct{}

func (SQLiteOpener) Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}
