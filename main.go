package main

import (
	"fmt"

	sqltranslation "github.com/vitalikir156/tasker/dbc"
	ui "github.com/vitalikir156/tasker/userinterface"
)

var killer bool

func main() {
	db := sqltranslation.Start("user=vit password=p!ssword2717 dbname=tasker port=5432 sslmode=disable")
	defer db.Close()
	for !killer {
		killer = ui.UI(db)
	}
	fmt.Println("Yeap! Exiting")
}
