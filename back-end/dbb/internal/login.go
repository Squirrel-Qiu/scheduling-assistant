package internal

import (
	"database/sql"
	"log"

	"golang.org/x/xerrors"
)

func (db *Impl) Login(openid string) (ok bool, err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return false, xerrors.Errorf("start transaction failed: %w", err)
	}

	err = tx.QueryRow("SELECT 1 FROM person WHERE openid=?", openid).Scan(new(int))

	switch {
	case err == nil:
		if err = tx.Commit(); err != nil {
			return false, xerrors.Errorf("commit transaction failed: %w", err)
		}
		return true, nil

	default:
		return false, xerrors.Errorf("scan failed: %w", err)

	case xerrors.Is(err, sql.ErrNoRows):
		log.Println("[DEBUG] err no row")
	}

	if _, err = tx.Exec("INSERT INTO person (openid) VALUES (?)", openid); err != nil {
		log.Printf("%+v", xerrors.Errorf("insert failed: %w", err))

		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("insert openid failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("commit transaction failed: %w", err)
	}

	return true, nil
}
