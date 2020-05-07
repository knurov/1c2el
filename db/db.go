package db

import (
	"context"
	"fmt"

	"knurov.ru/el/1c2el/helper"
)

//Transformer add / update transformer
func Transformer(hlp *helper.Helper, fullName string, transType string) (id uint32) {
	hlp.Log.Debug(fmt.Sprintf("on transformer %v %v", fullName, transType))
	connect, err := hlp.Conf.Database.Pool.Acquire(context.Background())
	defer connect.Release()
	hlp.Log.Fatal("On acquire pool connection %v", err)

	rows, err := connect.Query(context.Background(),
		"select id from transformer where fullName = $1 and type = $2",
		fullName, transType)
	defer rows.Close()
	hlp.Log.Error("On check transformer - %v", err)

	// fmt.Println(rows.RawValues())
	for rows.Next() {
		rows.Scan(&id)
		break
	}

	if id != 0 {
		return id
	}

	row := connect.QueryRow(context.Background(),
		"insert into transformer (fullName, type) values($1, $2) RETURNING id ",
		fullName, transType)

	hlp.Log.Error(row.Scan(&id))
	return id
}

// func db(hlp *helper.Helper) {
// 	// postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]

// 	// postgresql://
// 	// postgresql://localhost
// 	// postgresql://localhost:5433
// 	// postgresql://localhost/mydb
// 	// postgresql://user@localhost
// 	// postgresql://user:secret@localhost
// 	// postgresql://other@localhost/otherdb?connect_timeout=10&application_name=myapp
// 	// postgresql://host1:123,host2:456/somedb?target_session_attrs=any&application_name=myapp

// 	connect, err := hlp.Conf.Database.Pool.Acquire(context.Background())
// 	defer connect.Release()
// 	hlp.Log.Fatal("On acquire pool connection %v", err)

// 	type Transformer struct {
// 		ID       int    `json:"id"`
// 		FullName string `json:"fullName"`
// 		Type     string `json:"type"`
// 	}

// 	trans := Transformer{FullName: "tlo1", Type: "type1"}

// 	row := connect.QueryRow(context.Background(),
// 		"insert into transformer (fullName, type) values($1, $2) RETURNING id ",
// 		trans.FullName, trans.Type)

// 	var id int
// 	hlp.Log.Error(row.Scan(&id))
// 	fmt.Println(id)
// }
