package dbb

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"schedule/dbb/internal"
	"schedule/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Interface interface {
	Login(openid string) (ok bool, err error)
	SavePerson(openid string, nickName string) (ok bool, err error)
	OpenidAndNickName(rotaId int64) (person map[string]string, err error)

	NewRota(rota model.Rota, openid string) (ok bool, err error)
	GetRotas(openid string) (rotas []model.Rota, err error)

	GetFrees(openid string, rotaId int64) (frees []int, err error)
	GetLimitChoose(rotaId int64) (limit int, err error)
	ChooseFree(openid string, rotaId int64, frees []int) (ok bool, err error)

	InitPerson(rotaId int64) (personShift map[string]int, err error)
	QueryRotaInfo(rotaId int64) (shift, counter int, err error)
	QueryFree(rotaId int64) (freeId []int, err error)
	QueryChoosePersons(rotaId int64, freeId int) (choosePerson []string, err error)

	DeleteRota(openid string, rotaId int64) (ok bool, err error)
	io.Closer
}

var DB Interface

func InitDB(user, password string) (err error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/schedule", user, password))

	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(200)
	db.SetConnMaxLifetime(300 * time.Second)

	DB = &internal.Impl{DB: db}

	return nil
}
