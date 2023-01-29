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
	var e reflect.Value = reflect.ValueOf(s)
	var sqlstatement string
	sqlstatement1 := "CREATE TABLE IF NOT EXISTS " + table + " ("
	for i := 0; i < e.NumField(); i++ {
		var vt string
		varName := e.Type().Field(i).Name
		sqlstatement += "," + varName + " "
		//sqlstatement += varName + " "
		varType := e.Type().Field(i).Type
		switch varType.Kind() {
		case reflect.Int:
			if varName == "ID" {
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
	// such a crappy way to do this
	sqlstatement = sqlstatement[1:]
	sqlstatement += ")"
	sqlstatement = sqlstatement1 + sqlstatement

	return sqlstatement
}

/*
func InsertStructSlice(table string, s []interface{}) string {
	var sqls string
	for _, v := range s {
		sqls += InsertIntoTable(table, v)
	}
	return sqls
}
*/

func InsertIntoTable(table string, s interface{}) string {

	//	fmt.Println("InsertIntoTable:", table)

	var middlesql1 string
	var middlesql2 string

	var e reflect.Value = reflect.ValueOf(s)

	middlesql1 = "INSERT INTO " + table + " ("
	middlesql2 = ")VALUES("
	for i := 0; i < e.NumField(); i++ {
		//var vt string
		varName := e.Type().Field(i).Name
		//sqlstatement += "," + varName + " "
		//sqlstatement += varName + " "
		varType := e.Type().Field(i).Type

		varValue := e.Field(i).Interface()

		middlesql1 += varName + ","
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

	/*
	   o := "BEGIN;\n"
	   beginstatement, err := database.Prepare(o)

	   	if err != nil {
	   		log.Panicln(err)
	   	}

	   _, err = beginstatement.Exec()

	   	if err != nil {
	   		log.Panicln(err)
	   	}

	   	for _, v := range Systems {
	   		//fmt.Println("Inserting system", ind)
	   		var ssql, vsql string
	   		e := reflect.ValueOf(&v).Elem()
	   		for i := 0; i < e.NumField(); i++ {
	   			varName := e.Type().Field(i).Name
	   			ssql += "," + varName
	   			varValue := e.Field(i).Interface()
	   			vv := fmt.Sprintf("%v", varValue)
	   			varType := e.Type().Field(i).Type
	   			switch varType.Name() {
	   			case "int":
	   				vsql += "," + vv
	   			case "string":
	   				vsql += "," + "'" + vv + "'"
	   			}
	   		}
	   		ssql = ssql[1:]
	   		vsql = vsql[1:]
	   		sqlstatement := "INSERT INTO system( " + ssql + " ) VALUES(" + vsql + ")"
	   		//fmt.Println(sqlstatement)
	   		statement, _ := database.Prepare(sqlstatement)
	   		statement.Exec()
	   	}

	   o = "COMMIT;\n"
	   commitstatement, err := database.Prepare(o)

	   	if err != nil {
	   		log.Panicln(err)
	   	}

	   _, err = commitstatement.Exec()

	   	if err != nil {
	   		log.Panicln(err)
	   	}
	*/
	middlesql1 = middlesql1[:len(middlesql1)-1]
	middlesql2 = middlesql2[:len(middlesql2)-1] + ");"
	yyy := middlesql1 + middlesql2
	return yyy
}
