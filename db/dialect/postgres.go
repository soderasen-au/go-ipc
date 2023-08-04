package dialect

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresOpener struct{}

func (PostgresOpener) Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
