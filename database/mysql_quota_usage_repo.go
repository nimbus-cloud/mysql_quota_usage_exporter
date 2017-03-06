package database

import (
	"database/sql"
	"fmt"
	"log"
)

const mysqlQuotaUsageQueryPattern = `
SELECT tables.table_schema AS db,
       ROUND(SUM(tables.data_length + tables.index_length) / 1024 / 1024, 1) as size_mb,
       MAX(instances.max_storage_mb) as max_storage_mb
FROM   information_schema.tables AS tables
JOIN   %s.service_instances AS instances ON tables.table_schema = instances.db_name COLLATE utf8_general_ci
GROUP  BY tables.table_schema
`

type Database struct {
	Name string
	SizeMB float64
	MaxStorageMB float64
}

type Repo interface {
	All() ([]Database, error)
}

func NewMysqlQuotaUsageRepo(brokerDBName string, db *sql.DB) Repo {
	query := fmt.Sprintf(mysqlQuotaUsageQueryPattern, brokerDBName)
	return newRepo(query, db)
}

type repo struct {
	query  string
	db     *sql.DB
}

func newRepo(query string, db *sql.DB) Repo {
	return &repo{
		query:  query,
		db:     db,
	}
}

func (r repo) All() ([]Database, error) {
	log.Println("Executing All")

	databases := []Database{}

	rows, err := r.db.Query(r.query)
	if err != nil {
		return databases, fmt.Errorf("Error executing All: %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var dbName string
		var sizeMB, maxStorageMB float64
		if err := rows.Scan(&dbName, &sizeMB, &maxStorageMB); err != nil {
			return databases, fmt.Errorf("Scanning result row for dbName, sizeMB, maxStorageMB: %s", err.Error())
		}

		databases = append(databases, Database{dbName, sizeMB, maxStorageMB})
	}
	if err := rows.Err(); err != nil {
		return databases, fmt.Errorf("Reading result row of All: %s", err.Error())
	}

	return databases, nil
}
