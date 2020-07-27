package internal

import (
	"golang.org/x/xerrors"
)

func (db *Impl) ChooseFree(openid string, rotaId int64, frees []int) (ok bool, err error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return false, xerrors.Errorf("start transaction failed: %w", err)
	}

	rows, err := tx.Query("SELECT free_id FROM free WHERE openid=? AND rota_id=?", openid, rotaId)
	switch {
	default:
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("select frees failed: %w", err)
	case err == nil:
	}
	defer rows.Close()

	// delete if rows already exists, or insert
	if rows.Next() {
		if _, err = tx.Exec("DELETE FROM free WHERE openid=? AND rota_id=?", openid, rotaId); err != nil {
			if err := tx.Rollback(); err != nil {
				return false, xerrors.Errorf("rollback transaction failed: %w", err)
			}
			return false, xerrors.Errorf("db delete chosen frees failed: %w", err)
		}
	}

	for _, freeId := range frees {
		result, err := tx.Exec("INSERT INTO free (openid, rota_id, free_id) VALUES (?, ?, ?)",
			openid, rotaId, freeId)

		if err != nil {
			if err := tx.Rollback(); err != nil {
				return false, xerrors.Errorf("rollback transaction failed: %w", err)
			}
			return false, xerrors.Errorf("insert chosen freeId failed: %w", err)
		}

		// unnecessary
		if affected, _ := result.RowsAffected(); affected != 1 {
			if err := tx.Rollback(); err != nil {
				return false, xerrors.Errorf("rollback transaction failed: %w", err)
			}
			return false, xerrors.New("insert chosen freeId failed")
		}
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return false, xerrors.Errorf("rollback transaction failed: %w", err)
		}
		return false, xerrors.Errorf("commit transaction failed: %w", err)
	}

	return true, nil
}
