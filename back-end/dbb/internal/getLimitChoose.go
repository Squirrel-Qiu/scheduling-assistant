package internal

import (
	"database/sql"
	"golang.org/x/xerrors"
)

func (db *Impl) GetLimitChoose(rotaId int64) (limitChoose int, err error) {
	err = db.DB.QueryRow("SELECT limit_choose FROM rota WHERE rota_id=?", rotaId).Scan(&limitChoose)

	switch {
	default:
		return -1, xerrors.Errorf("select limit_choose failed: %w", err)

	case xerrors.Is(err, sql.ErrNoRows):
		return -1, err

	case err == nil:
	}

	return limitChoose, nil
}
