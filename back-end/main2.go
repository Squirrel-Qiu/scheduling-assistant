package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	user := "root"
	password := "root"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/schedule", user, password))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rotaId := 291255583271555078

	const sqlStr1 = "SELECT DISTINCT openid FROM free WHERE rota_id=?"

	rows1, err := db.Query(sqlStr1, rotaId)
	if err != nil {
		fmt.Printf("select openid failed: %v", err)
	}
	defer rows1.Close()

	var openids []string
	for rows1.Next() {
		var openid string
		if err := rows1.Scan(&openid); err != nil {
			fmt.Printf("scan failed: %v", err)
		}
		openids = append(openids, openid)
	}

	// create In(?,?...?,?)
	s := "?"
	for i := 1; i < len(openids); i++ {
		s += ",?"
	}

	// interface slice
	openidsInterface := make([]interface{}, len(openids))
	for i, v := range openids {
		openidsInterface[i] = v
	}

	sqlStr2 := fmt.Sprintf("SELECT openid, nick_name FROM person WHERE openid IN (%s)", s)
	rows2, err := db.Query(sqlStr2, openidsInterface...)
	if err != nil {
		fmt.Printf("select openid and nick_name failed: %v", err)
	}
	defer rows2.Close()

	person := make(map[string]string)
	for rows2.Next() {
		var o, n string
		if err := rows2.Scan(&o, &n); err != nil {
			fmt.Printf("scan failed: %v", err)
		}
		person[o] = n
	}
}
