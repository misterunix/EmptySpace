/*
Functions to create and manage database tables
*/
package main

import (
	"fmt"
	"reflect"
)

// Create table based on struct.
// Retuns the sql statement as a string.
// This is a work in progress.
func CreateTableFromStruct(table string, s interface{}) string {

	var reflectedValue reflect.Value = reflect.ValueOf(s) // reflect the struct (interface)

	var sqlstatement string

	sqlstatement1 := "CREATE TABLE IF NOT EXISTS " + table + " ("
	for i := 0; i < reflectedValue.NumField(); i++ {
		var vt string
		varName := reflectedValue.Type().Field(i).Name // get the name of the field
		sqlstatement += "," + varName + " "
		varType := reflectedValue.Type().Field(i).Type // get the type of the field

		// Did this differnt than the other reflect code. This is a work in progress.
		switch varType.Kind() {
		case reflect.Int:
			if varName == "ID" { // detect if the field is the ID field
				vt = "INTEGER NOT NULL PRIMARY KEY"
			} else {
				vt = "INTEGER"
			}
		case reflect.String:
			vt = "TEXT"
		case reflect.Float64:
			vt = "REAL"
		case reflect.Bool:
			vt = "INTEGER"
		}
		sqlstatement += vt
	}

	// such a crappy way to do this. Return to this at a later date.
	sqlstatement = sqlstatement[1:] // remove the first comma
	sqlstatement += ")"
	sqlstatement = sqlstatement1 + sqlstatement

	return sqlstatement
}

func InsertIntoTable(table string, s interface{}) string {

	var middlesql1 string
	var middlesql2 string

	var reflectedValue reflect.Value = reflect.ValueOf(s)

	middlesql1 = "INSERT INTO " + table + " ("
	middlesql2 = ")VALUES("
	for i := 0; i < reflectedValue.NumField(); i++ {

		varName := reflectedValue.Type().Field(i).Name
		varType := reflectedValue.Type().Field(i).Type
		varValue := reflectedValue.Field(i).Interface()

		middlesql1 += varName + ","

		// This is my normal way of working with reflect. Strings may be slower but easier to read.
		switch varType.Name() {
		case "int":
			middlesql2 += fmt.Sprintf("%d", varValue.(int)) + ","
		case "string":
			middlesql2 += "'" + varValue.(string) + "',"
		case "float64":
			middlesql2 += fmt.Sprintf("%f", varValue.(float64)) + ","
		case "bool":
			middlesql2 += fmt.Sprintf("%v", varValue.(bool)) + ","
		default:
			return ""
		}
	}

	middlesql1 = middlesql1[:len(middlesql1)-1]
	middlesql2 = middlesql2[:len(middlesql2)-1] + ");"
	yyy := middlesql1 + middlesql2
	return yyy
}
