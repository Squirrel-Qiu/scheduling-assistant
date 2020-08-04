package internal

import "golang.org/x/xerrors"

func (db *Impl) SavePerson(openid, nickName string) (ok bool, err error) {
	sm, err := db.DB.Prepare("INSERT INTO person (openid,nick_name) VALUES (?,?) " +
		"on duplicate key UPDATE nick_name=VALUES(nick_name)")
	if err != nil {
		return false, xerrors.Errorf("prepare sm failed: %w", err)
	}

	_, err = sm.Exec(openid, nickName)
	if err != nil {
		return false, xerrors.Errorf("sm exec failed: %w", err)
	}

	return true, nil
}
