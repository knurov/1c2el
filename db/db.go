package db

import (
	"context"
	"fmt"

	"knurov.ru/el/1c2el/helper"
)

//applyMap inser values
func applyMap(hlp *helper.Helper, selectStatment string, insertStatment string, order []string, values map[string]string) (id uint32) {
	connect, err := hlp.Conf.Database.Pool.Acquire(context.Background())
	defer connect.Release()
	hlp.Log.Fatal("On acquire pool connection %v", err)

	if len(values) != len(order) {
		hlp.Log.Fatal("Count fuields and values not equal")
	}

	paramsList := make([]interface{}, 0, len(order))
	for _, fieldName := range order {
		paramsList = append(paramsList, values[fieldName])
	}

	rows, err := connect.Query(context.Background(), selectStatment, paramsList...)
	defer rows.Close()
	hlp.Log.Error("On select - %v", err)

	for rows.Next() {
		rows.Scan(&id)
		hlp.Log.Trace("Found record %v", paramsList)
		return id
	}

	row := connect.QueryRow(context.Background(), insertStatment, paramsList...)

	hlp.Log.Error(row.Scan(&id))
	return id

}

//Transformer add/update transformer
func Transformer(hlp *helper.Helper, values map[string]string) (id uint32) {
	hlp.Log.Debug(fmt.Sprintf("on transformer %#v", values))
	selectStatment := "select id from transformer where fullName = $1 and type = $2"
	insertStatment := "insert into transformer (fullName, type) values($1, $2) RETURNING id "
	order := []string{"fullName", "type"}
	return applyMap(hlp, selectStatment, insertStatment, order, values)

}

//transformer add/update transformer
func transformer(hlp *helper.Helper, fullName string, transType string) (id uint32) {
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
