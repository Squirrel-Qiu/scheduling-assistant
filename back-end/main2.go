package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"schedule/model"
)

func main() {
	user := "root"
	password := "root"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/schedule", user, password))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id := 201703
	//rotaId := 2

	const sqlStr = "SELECT rota_id, title, shift, limit_choose, counter FROM rota WHERE id=?"
	//err = db.QueryRow(sqlStr, id).Scan(1)
	//if xerrors.Is(err, sql.ErrNoRows) {
	//	fmt.Println(xerrors.Errorf("the error: %w", err))
	//	return
	//}
	//fmt.Println(err)

	rows, err := db.Query(sqlStr, id)
	if err != nil {
		fmt.Printf("query failed: %v", err)
	}
	defer rows.Close()

	var rotas []model.Rota
	for rows.Next() {
		var rota model.Rota
		if err := rows.Scan(&rota.RotaId, &rota.Title, &rota.Shift, &rota.LimitChoose, &rota.Counter); err != nil {
			fmt.Printf("scan failed: %v", err)
		}
		rotas = append(rotas, rota)
	}

	fmt.Println(rotas)
}
