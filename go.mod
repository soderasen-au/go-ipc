module github.com/soderasen-au/go-ipc

go 1.20

require (
	github.com/rs/zerolog v1.29.1
	github.com/soderasen-au/go-common v0.2.0
	github.com/xhit/go-simple-mail/v2 v2.15.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/driver/sqlite v1.5.2
	gorm.io/driver/sqlserver v1.5.1
	gorm.io/gorm v1.25.2
)

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/go-test/deep v1.1.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/microsoft/go-mssqldb v1.5.0 // indirect
	github.com/toorop/go-dkim v0.0.0-20201103131630-e1cd1a0a5208 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
)

replace github.com/soderasen-au/go-common v0.2.0 => ../../soderasen-au/go-common
