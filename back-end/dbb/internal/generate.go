package internal

import (
	"database/sql"
	"golang.org/x/xerrors"
)

func (db *Impl) InitPerson(rotaId int64) (personShift map[string]int, err error) {
	rows, err := db.DB.Query("SELECT DISTINCT openid from free where rota_id=?", rotaId)
	if err != nil {
		return nil, xerrors.Errorf("select persons who have choose failed: %w", err)
	}
	defer rows.Close()

	personShift = make(map[string]int)
	for rows.Next() {
		var person string
		if err = rows.Scan(&person); err != nil {
			return nil, xerrors.Errorf("scan personShift's person failed: %w", err)
		}

		personShift[person] = 0
	}

	return personShift, nil
}

func (db *Impl) QueryRotaInfo(rotaId int64) (shift, counter int, err error) {
	err = db.DB.QueryRow("SELECT shift,counter FROM rota WHERE rota_id=?", rotaId).Scan(&shift, &counter)

	switch {
	default:
		return -1, -1, xerrors.Errorf("select rota's info(shift counter) failed: %w", err)

	case xerrors.Is(err, sql.ErrNoRows):
		return -1, -1, err

	case err == nil:
	}

	return shift, counter, nil
}

// 已填的所有时间段： 最少人选择的空闲时间段 ～ 最多人选择的空闲时间段
func (db *Impl) QueryFree(rotaId int64) (frees []int, err error) {
	rows, err := db.DB.Query("SELECT free_id from free where rota_id=? group by free_id order by count(*)", rotaId)
	if err != nil {
		return nil, xerrors.Errorf("select frees order failed: %w", err)
	}
	defer rows.Close()

	var free int
	for rows.Next() {
		if err = rows.Scan(&free); err != nil {
			return nil, xerrors.Errorf("scan frees failed: %w", err)
		}
		frees = append(frees, free)
	}

	return frees, nil
}

// 选择该时间段的所有人
func (db *Impl) QueryChoosePersons(rotaId int64, freeId int) (choosePersons []string, err error) {
	rows, err := db.DB.Query("SELECT openid FROM free WHERE rota_id=? and free_id=?", rotaId, freeId)
	if err != nil {
		return nil, xerrors.Errorf("select the persons of this freeId failed: %w", err)
	}
	defer rows.Close()

	var person string
	for rows.Next() {
		if err = rows.Scan(&person); err != nil {
			return nil, xerrors.Errorf("scan choose persons failed: %w", err)
		}
		choosePersons = append(choosePersons, person)
	}

	return choosePersons, nil
}
