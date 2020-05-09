package db

import (
	"context"
	"fmt"

	"knurov.ru/el/1c2el/helper"
)

//Transformer add/update transformer
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

	for rows.Next() {
		rows.Scan(&id)
		hlp.Log.Trace("Found transformer %v", fullName)
		return id
	}

	hlp.Log.Trace("Insert new transformer %v", fullName)
	row := connect.QueryRow(context.Background(),
		"insert into transformer (fullName, type) values($1, $2) RETURNING id ",
		fullName, transType)

	hlp.Log.Error(row.Scan(&id))
	return id
}
