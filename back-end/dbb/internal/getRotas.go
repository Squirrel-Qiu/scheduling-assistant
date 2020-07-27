package internal

import (
	"database/sql"
	"golang.org/x/xerrors"
	"schedule/model"
)

func (db *Impl) GetRotas(openid string) (rotas []model.Rota, err error) {
	const sqlStr = "SELECT rota_id, title, shift, limit_choose, counter FROM rota WHERE openid=?"

	err = db.DB.QueryRow(sqlStr, openid).Scan(new(int))
	if xerrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	rows, err := db.DB.Query(sqlStr, openid)
	if err != nil {
		return nil, xerrors.Errorf("select rotas failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rota model.Rota
		if err := rows.Scan(&rota.RotaId, &rota.Title, &rota.Shift, &rota.LimitChoose, &rota.Counter); err != nil {
			return nil, xerrors.Errorf("scan rotas info failed: %w", err)
		}
		rotas = append(rotas, rota)
	}

	return rotas, nil
}
