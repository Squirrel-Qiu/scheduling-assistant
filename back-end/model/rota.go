package model

type Rota struct {
	RotaId       int64  `json:"rota_id,string"`
	Title	     string  `json:"title"`
	Shift	     int     `json:"shift"`
	LimitChoose  int	 `json:"limit_choose"`
	Counter      int     `json:"counter"`
}

type SimpleRota struct {
	RotaId       int64  `json:"rota_id,string"`
	Title	     string  `json:"title"`
}
