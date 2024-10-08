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

func insert(db *sqlx.DB) {
	var task sqltranslation.Tasktable
	fmt.Printf("enter message:")
	var err error
	task.Message, err = userinputreader()
	if err != nil {
		fmt.Println(err)
	}
	if len(task.Message) < 1 {
		fmt.Println("No task message entered, exiting")
	}

	fmt.Printf("enter status(empty for NEW):")
	task.Status, err = userinputreader()
	if err != nil {
		fmt.Println(err)
	}
	if len(task.Status) < 1 {
		task.Status = "new"
	}
	task.Deadline, err = datetimegetter()
	if err != nil {
		fmt.Println(err)
		fmt.Println("task without deadline is fiasco, exiting")
		return
	}
	fmt.Printf("status:%v, deadline:%v, message:%v\n", task.Status, task.Deadline, task.Message)
	res, err := sqltranslation.PostTask(db, task)

	fmt.Println(res, err)
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
				date, err2 := datetimegetter()
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

func get(db *sqlx.DB) {
	res, err := sqltranslation.GetAll(db)
	fmt.Println(res, err)
}

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

func userinputreader() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	text, err := reader.ReadString('\n')
	text = strings.ReplaceAll(text, "\n", "")
	// fmt.Println(err)
	return text, err
}

func datetimegetter() (time.Time, error) {
	nowdate := time.Now().Format("02.01.06")
	nowtime := time.Now().Format("15.04")
	fmt.Printf("enter deadline date (format dd.mm.yy)(nothing for %v)", nowdate)
	dateread, err2 := userinputreader()
	if err2 != nil {
		fmt.Println(err2)
	}
	if len(dateread) > 0 {
		nowdate = dateread
	}
	fmt.Printf("enter deadline time (format hh.mm)(nothing for %v)", nowtime)
	timeread, err2 := userinputreader()
	if err2 != nil {
		fmt.Println(err2)
	}
	if len(timeread) > 0 {
		nowtime = timeread
	}
	date, err := time.Parse("02.01.0615.04", nowdate+nowtime)
	if err != nil {
		fmt.Println(err)
	}
	return date, err
}
