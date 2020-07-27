package internal

import (
	"database/sql"
	"golang.org/x/xerrors"
)

func (db *Impl) GetFrees(openid string, rotaId int64) (frees []int, err error) {
	const sqlStr = "SELECT free_id FROM free WHERE openid=? AND rota_id=?"

	err = db.DB.QueryRow(sqlStr, openid).Scan(new(int))
	if xerrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	rows, err := db.DB.Query(sqlStr, openid, rotaId)
	if err != nil {
		return nil, xerrors.Errorf("select frees failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var free int
		if err := rows.Scan(&free); err != nil {
			return nil, xerrors.Errorf("scan frees failed: %w", err)
		}
		frees = append(frees, free)
	}

	return frees, nil
}
