package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	sqltranslation "github.com/vitalikir156/tasker/dbc"
)

type Switchertype struct {
	Userinput      string
	Menu0          int
	Stage          int
	Allowuserinput bool
}

func Post(menu *Switchertype, task *sqltranslation.Tasktable, db *sqlx.DB) (string, error) {
	var err error
	switch menu.Stage {
	case 0:
		{
			menu.Stage = +1
			return "enter task message:", nil
		}
	case 1:
		{
			if len(menu.Userinput) < 1 {
				return "Empty message entered, please try again or exit", nil
			} else {
				task.Message = menu.Userinput
				menu.Stage = +1
				return "enter status(empty for <new>):", nil
			}
		}
	case 2:
		{
			if len(menu.Userinput) < 1 {
				task.Status = "new"
			} else {
				task.Status = menu.Userinput
				menu.Stage = +1
			}
			return "enter status(empty for <new>):", nil
		}
	case 3:
		{
		}
	}
	task.Deadline, err = Datetimegetter()
	if err != nil {
		fmt.Println(err)
		fmt.Println("task without deadline is fiasco, exiting")
	}
	fmt.Printf("status:%v, deadline:%v, message:%v\n", task.Status, task.Deadline, task.Message)
	res, err := sqltranslation.PostTask(db, *task)

	fmt.Println(res, err)
	return "enter status(empty for <new>):", nil
}

func update(db *sqlx.DB) {
	var killer bool
	var task sqltranslation.Tasktable
	fmt.Println("enter ID for update:")
	_, err := fmt.Scanln(&task.ID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("no message ID provided. Exiting")
	}
	res, err := sqltranslation.GetOverID(db, task)
	if err != nil {
		fmt.Println(err)
	}
	if len(res) < 1 {
		fmt.Printf("task with ID %v not found. Exiting\n", task.ID)
		return
	}
	for !killer {
		fmt.Println(res)
		fmt.Print(`-=UPDATE TASK=-
	0: save and exit to main menu
	1: exit without saving
	2: update message
	3: update status
	4: update deadline
	select option and hit enter:`)
		var s int
		_, err = fmt.Scanln(&s)
		if err != nil {
			fmt.Println(err)
			s = -1
		}
		switch s {
		case 0:
			{
				killer = true
				out, err2 := sqltranslation.EditTask(db, task)
				if err != nil {
					fmt.Println(err2)
				} else {
					fmt.Println(out)
				}
			}
		case 1:
			{
				killer = true
			}
		case 2:
			{
				fmt.Printf("enter message:")
				task.Message, err = userinputreader()
				if err != nil {
					fmt.Println(err)
				}
				if len(task.Message) < 1 {
					fmt.Println("No task message entered, exiting")
				}
			}
		case 3:
			{
				fmt.Printf("enter status(empty for NEW):")
				task.Status, err = userinputreader()
				if err != nil {
					fmt.Println(err)
				}
				if len(task.Status) < 1 {
					fmt.Println("No status entered, exiting")
				}
			}

		case 4:
			{
				date, err2 := Datetimegetter()
				if err2 == nil {
					task.Deadline = date
				}
			}
		default:
			{
				fmt.Printf("bad input entered: %v\n", s)
			}
		}
	}
}

func del(db *sqlx.DB) {
	var task sqltranslation.Tasktable
	fmt.Println("enter ID for delete:")
	_, err := fmt.Scanln(&task.ID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("no message ID provided. Exiting")
	}
	res, err := sqltranslation.DelTask(db, task)

	fmt.Println(res, err)
}

/*
	func UI(db *sqlx.DB) bool {
		fmt.Print(`
		-=MAIN MENU=-

0: exit
1: get all tasks
2: create task
3: edit task
4: delete task

		select option and hit enter:`)
		var s int
		_, err := fmt.Scanln(&s)
		if err != nil {
			fmt.Println(err)
			return false
		}
		switch s {
		case 0:
			{
				return true
			}
		case 1:
			{
				get(db)
			}
		case 2:
			{
				insert(db)
			}
		case 3:
			{
				update(db)
			}
		case 4:
			{
				del(db)
			}
		default:
			{
				fmt.Printf("bad input entered: %v\n", s)
			}
		}
		return false
	}
*/
func Mainmenu(menu *Switchertype) (string, error) {

	switch menu.Userinput {
	case "":
		{
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
	case "0":
		{
			menu.Menu0 = -1
			menu.Allowuserinput = false
			return "exiting", nil
		}
	case "1":
		{
			menu.Menu0 = 1
			menu.Allowuserinput = false
			return "", nil
		}
	case "2":
		{
			menu.Menu0 = 2
			//		menu.Allowuserinput = false
			return "", nil
		}
	}
	menu.Allowuserinput = false
	return fmt.Sprintf("Something very strange was entered: %v.\nEnter an existing option.", menu.Userinput), nil
}
func userinputreader() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	text, err := reader.ReadString('\n')
	text = strings.ReplaceAll(text, "\n", "")
	// fmt.Println(err)
	return text, err
}
func UpdateDate(menu *Switchertype, task *sqltranslation.Tasktable) (string, error) {
	fmt.Println(len(menu.Userinput))
	if len(menu.Userinput) == 0 {
		return fmt.Sprintf("entered nothing, keeping date=%v", task.Deadline.Format("02.01.06")), nil
	}
	if len(menu.Userinput) != 8 {
		return "", fmt.Errorf("entered bad sequence: %v, need dd.mm.yy as 02.01.06", menu.Userinput)

	}
	now, err := time.Parse("02.01.0615.04", menu.Userinput+task.Deadline.Format("15.04"))
	if err != nil {
		return "", err
	}
	task.Deadline = now
	return "date updated without any errors", nil
}

func Datetimegetter() (time.Time, error) {
	var err error
	now := time.Now()
	fmt.Println(now)
	//.Format("02.01.06")
	//	nowtime := time.Now().Format("15.04")
	fmt.Printf("enter deadline date (format dd.mm.yy)(nothing for %v)", now.Format("02.01.06"))
	dateread, err2 := userinputreader()
	if err2 != nil {
		fmt.Println(err2)
	}
	if len(dateread) > 0 {
		now, err = time.Parse("02.01.0615.04", dateread+now.Format("15.04"))
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(now)
	fmt.Printf("enter deadline time (format hh.mm)(nothing for %v)", now.Format("15.04"))
	timeread, err2 := userinputreader()
	if err2 != nil {
		fmt.Println(err2)
	}

	if len(dateread) > 0 {
		now, err = time.Parse("02.01.0615.04", now.Format("02.01.06")+timeread)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(now)

	return now, err
}
