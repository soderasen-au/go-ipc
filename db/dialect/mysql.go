package dialect

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlOpener struct{}

func (MysqlOpener) Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
