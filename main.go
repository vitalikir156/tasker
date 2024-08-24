package main

import (
	"fmt"
	"time"

	sqltranslation "github.com/vitalikir156/tasker/dbc"
)

func srv(url string) {
	db := sqltranslation.Start("user=vit password=p!ssword2717 dbname=market host=db port=5432 sslmode=disable")
	defer db.Close()

}

func main() {
	db := sqltranslation.Start("user=vit password=p!ssword2717 dbname=tasker port=5432 sslmode=disable")
	timer, err := time.Parse("02.01.06", "17.03.98")
	fmt.Println(timer)

	task := sqltranslation.Tasker{
		Message:  "Create your own program3!",
		Deadline: timer,
	}

	res, err := sqltranslation.PostTask(db, task)
	defer db.Close()
	fmt.Println(res, err)
	fmt.Println("Yeap! Exiting")
}
