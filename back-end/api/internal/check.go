package internal

import "schedule/model"

func newRotaCheck(rota *model.Rota) bool {
	// check shift:34 limit_choose:35 counter:1-100
	if rota.Shift < 1 || rota.Shift > 34 {
		return false
	}

	if rota.LimitChoose < rota.Shift || rota.LimitChoose > 35 {
		return false
	}

	if rota.Counter < 1 {
		return false
	}

	return true
}
