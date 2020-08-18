package internal

import (
	"fmt"
	"golang.org/x/xerrors"
	"schedule/model"
)

func (db *Impl) GetJoin(openid string) (joins []model.SimpleRota, err error) {
	rows, err := db.DB.Query("SELECT DISTINCT rota_id FROM free WHERE openid=?", openid)
	if err != nil {
		return nil, xerrors.Errorf("select rotaId who have join failed: %w", err)
	}
	defer rows.Close()

	var rotas []int64
	for rows.Next() {
		var rotaId int64
		if err = rows.Scan(&rotaId); err != nil {
			return nil, xerrors.Errorf("scan join rota failed: %w", err)
		}
		rotas = append(rotas, rotaId)
	}

	// create In(?,?...?,?)
	s := "?"
	for i := 1; i < len(rotas); i++ {
		s += ",?"
	}

	// rotas []int64 To []interface{}
	rotasInterface := make([]interface{}, len(rotas))
	for i, v := range rotas {
		rotasInterface[i] = v
	}

	sqlStr := fmt.Sprintf("SELECT rota_id, title FROM rota WHERE rota_id IN (%s)", s)

	rows2, err := db.DB.Query(sqlStr, rotasInterface...)
	if err != nil {
		return nil, xerrors.Errorf("select rota_id and title failed: %w", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var join model.SimpleRota
		if err = rows2.Scan(&join.RotaId, &join.Title); err != nil {
			return nil, xerrors.Errorf("scan join info failed: %w", err)
		}
		joins = append(joins, join)
	}

	return joins, nil
}


