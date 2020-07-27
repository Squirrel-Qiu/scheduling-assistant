package internal

import "golang.org/x/xerrors"

func (db *Impl) DeleteRota(openid string, rotaId int64) (ok bool, err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return false, xerrors.Errorf("start transaction failed: %w", err)
	}

	if _, err = tx.Exec("DELETE FROM free WHERE rota_id=?", rotaId); err != nil {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("db delete rota's frees failed: %w", err)
	}

	result, err := tx.Exec("DELETE FROM rota WHERE rota_id=? AND openid=?", rotaId, openid)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("db delete rota failed: %w", err)
	}

	// unnecessary
	if affected, _ := result.RowsAffected(); affected != 1 {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.New("db delete rota failed, delete result is not one")
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("commit transaction failed: %w", err)
	}

	return true, nil
}
