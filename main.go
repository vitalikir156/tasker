package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	sqltranslation "github.com/vitalikir156/tasker/dbc"
	ui "github.com/vitalikir156/tasker/userinterface"
)

func main() {
	var task sqltranslation.Tasktable
	var switcher ui.Switchertype
	switcher.Allowuserinput = true
	db := sqltranslation.Start("user=vit password=p!ssword2717 dbname=tasker port=5432 sslmode=disable")
	defer db.Close()

	for switcher.Menu0 > -1 {
		var err error

		out, err := menucore(&switcher, &task, db)
		if err != nil {
			fmt.Println(err)
			err = nil
		} else {
			fmt.Println(out)
		}
		if switcher.Allowuserinput {
			switcher.Userinput, err = userinputreader()
		} else {
			switcher.Allowuserinput = true
			switcher.Userinput = ""
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	// var selector
}

func menucore(menu *ui.Switchertype, task *sqltranslation.Tasktable, db *sqlx.DB) (string, error) {
	switch menu.Userinput {
	case "exit": // Выход из утилиты в любом месте
		{
			menu.Allowuserinput = false
			menu.Menu0 = -1
			return "exiting", nil
		}
	case "main": // Выход в основное меню в любом месте
		{
			menu.Userinput = ""
			menu.Menu0 = 0
			task = &sqltranslation.Tasktable{}
		}
	}
	switch menu.Menu0 {
	case 0:
		{
			out, err := ui.Mainmenu(menu)

			return out, err
		}

	case 1:

		{
			menu.Allowuserinput = false
			menu.Menu0 = 0
			return sqltranslation.GetAll(db)
		}
	case 2:
		{
			out, err := ui.Post(menu, task, db)
			return out, err
		}
	case 3:
		{
			out, err := ui.Update(menu, task, db)
			return out, err
		}
	case 4:
		{
			out, err := ui.Del(menu, task, db)
			return out, err
		}
	case 11:
		{
			out, err := ui.GetWithStatus(menu, db)
			return out, err
		}
	case 12:
		{
			out, err := ui.GetWithoutStatus(menu, db)
			return out, err
		}
	case 13:
		{
			menu.Allowuserinput = false
			menu.Menu0 = 0
			return sqltranslation.GetAllDeadlined(db)
		}
	}
	return "assss", nil
}

func userinputreader() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	text, err := reader.ReadString('\n')
	text = strings.ReplaceAll(text, "\n", "")
	// fmt.Println(err)
	return text, err
}
