package internal

import (
	"golang.org/x/xerrors"
	"schedule/model"
	"time"
)

func (db *Impl) NewRota(rota model.Rota, openid string) (ok bool, err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return false, xerrors.Errorf("start transaction failed: %w", err)
	}

	result, err := tx.Exec("INSERT INTO rota (rota_id, title, openid, shift, limit_choose, counter, date) VALUES (?, ?, ?, ?, ?, ?, ?)",
		rota.RotaId, rota.Title, openid, rota.Shift, rota.LimitChoose, rota.Counter, time.Now())

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("insert new rota failed: %w", err)
	}

	if affected, _ := result.RowsAffected(); affected != 1 {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.New("insert new rota failed")
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("commit transaction failed: %w", err)
	}

	return true, nil
}
