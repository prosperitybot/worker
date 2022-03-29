package internal

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Database *sqlx.DB

func OpenDatabase(connectionString string) error {
	var err error

	Database, err = sqlx.Connect("mysql", connectionString)

	if err != nil {
		return err
	}

	return nil
}

func DbConnectionString(host string, user string, pass string, db string, port string, sslmode string, timezone string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, pass, db, port, sslmode, timezone)
}
