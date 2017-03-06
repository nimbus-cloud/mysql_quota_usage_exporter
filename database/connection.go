package database

import (
	"database/sql"
	"fmt"
)

func NewConnection(user, password, host, dbName string, port int) (*sql.DB, error) {

	var userPass string
	if password != "" {
		userPass = fmt.Sprintf("%s:%s", user, password)
	} else {
		userPass = user
	}

	return sql.Open("mysql", fmt.Sprintf(
		"%s@tcp(%s:%d)/%s",
		userPass,
		host,
		port,
		dbName,
	))
}
