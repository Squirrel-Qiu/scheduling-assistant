package internal

import (
	"database/sql"
	"golang.org/x/xerrors"
)

func (db *Impl) Login(openid string) (ok bool, err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return false, xerrors.Errorf("start transaction failed: %w", err)
	}

	err = tx.QueryRow("SELECT 1 FROM person WHERE openid=?", openid).Scan(new(int))

	switch {
	default:
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		if err != nil {
			return false, xerrors.Errorf("check person exists failed: %w", err)
		}
		return true, nil
	case xerrors.Is(err, sql.ErrNoRows):
	}

	if _, err = tx.Exec("INSERT INTO person (openid) VALUES (?)", openid); err != nil {
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
