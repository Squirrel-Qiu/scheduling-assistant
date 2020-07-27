package internal

import "database/sql"

type Impl struct {
	DB *sql.DB
}

func (db *Impl) Close() error {
	return db.DB.Close()
}
