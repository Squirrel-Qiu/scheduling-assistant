package snowid

import (
	"log"

	"github.com/holdno/snowFlakeByGo"
	"golang.org/x/xerrors"
)

type RotaidGetter interface {
	GetRotaId() int64
}

type SnowFlake struct{}

func (sf SnowFlake) GetRotaId() int64 {
	worker, err := snowFlakeByGo.NewWorker(0)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("new worker failed: %w", err))
	}

	return worker.GetId()
}
