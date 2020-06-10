package db

import (
	"context"

	"knurov.ru/el/1c2el/helper"
)

//applyMap inser values
func execQueryByMap(hlp *helper.Helper, selectStatment string, insertStatment string, order []string, values map[string]string) (id uint32) {
	if hlp.Conf.Database.DryRun {
		hlp.Log.Trace("Dry run mode! Skip of persisting %#q", values)
		return
	}
	hlp.Log.Trace("Persisting of values %#q", values)
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

//SerialNumber add/update SerialNumber
func SerialNumber(hlp *helper.Helper, values map[string]string) (id uint32) {
	hlp.Log.Debug("on persisting serial number %q", values["SerialNumber"])
	selectStatment := "select id from sn where fullName = $1 and type = $2"
	insertStatment := "insert into sn (fullName, type) values($1, $2) RETURNING id "
	order := []string{"fullName", "type"}
	return execQueryByMap(hlp, selectStatment, insertStatment, order, values)
}

//Transformer add/update transformer
func Transformer(hlp *helper.Helper, values map[string]string) (id uint32) {
	hlp.Log.Debug("on persisting Transformer %q", values["FullName"])
	selectStatment := "select id from transformer where FullName = $1 and type = $2"
	insertStatment := "insert into transformer (fullName, type) values($1, $2) RETURNING id "
	order := []string{"fullName", "type"}
	return execQueryByMap(hlp, selectStatment, insertStatment, order, values)
}

//Coil add/update transformer
func Coil(hlp *helper.Helper, values map[string]string) (id uint32) {
	hlp.Log.Debug("on persisting coil %q", values["FullName"])
	selectStatment := "select id from transformer where fullName = $1 and type = $2"
	insertStatment := "insert into transformer (fullName, type) values($1, $2) RETURNING id "
	order := []string{"fullName", "type"}
	return execQueryByMap(hlp, selectStatment, insertStatment, order, values)
}
