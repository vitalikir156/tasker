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
	//var selector
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
			task = &sqltranslation.Tasktable{}
			fmt.Println(task)
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
			out, err := ui.UpdateDate(menu, task)

			return out, err
		}
	case 3:
		{
			out := "3"
			return out, nil
		}
	case 4:
		{
			out := "1"
			return out, nil
		}
	case 5:
		{
			out := "1"
			return out, nil
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

/*	if menu.Menu0 == 0 {

		{
			if len(menu.Userinput) > 0 {
				x, _ := strconv.Atoi(menu.Userinput)
				if 0 <= x && x <= 5 {
					menu.Menu0 = x
				}
			} else {
				a := (`
-=MAIN MENU=-
0: exit
1: get all tasks
2: create task
3: edit task
4: delete task
select option and hit enter:`)
				return a, nil
			}
		}
	}

	if menu.Menu0 == 1 {
		{
			a := (`

1: get all tasks
`)
			return a, nil
		}
	}
	menu.Allowuserinput = false*/
