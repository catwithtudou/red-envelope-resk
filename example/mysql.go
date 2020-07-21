package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	dsName := "root:@tcp(127.0.0.1:3306)/demo?charset=utf8&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		fmt.Println(err)
	}
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(3)
	db.SetConnMaxLifetime(7 * time.Hour)

	fmt.Println(db.Query("select now() "))

	defer db.Close()
}
