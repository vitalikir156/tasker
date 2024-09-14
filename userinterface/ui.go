package ui

import (
	"fmt"
	"strconv"
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
	switch menu.Stage {
	case 0:
		{
			menu.Stage++
			return "enter task message:", nil
		}
	case 1:
		{
			if len(menu.Userinput) < 1 {
				return "Empty message entered, please try again or exit", nil
			}
			task.Message = menu.Userinput
			menu.Stage++
			return "enter status(empty for <new>):", nil
		}
	case 2:
		{
			menu.Stage++
			if len(menu.Userinput) < 1 {
				task.Status = "new"
			} else {
				task.Status = menu.Userinput
			}
			task.Deadline = time.Now()
			return fmt.Sprintf("enter time of deadline (enter for %v)", task.Deadline.Format("15.04")), nil
		}
	case 3:
		{
			str, err := UpdateTime(menu, task)
			if err != nil {
				return "", err
			}
			menu.Stage++
			return fmt.Sprintf("%v\nenter date of deadline (enter for %v)", str, task.Deadline.Format("02.01.06")), err
		}
	case 4:
		{
			str, err := UpdateDate(menu, task)
			if err != nil {
				return "", err
			}
			menu.Stage++
			blablabla := "Y for save task, E for exit without saving"
			return fmt.Sprintf("%v\n Check task:\n status:%v, deadline:%v, message:%v\n%v\n ",
				str, task.Status, task.Deadline.Format("02.01.06 15.04"), task.Message, blablabla), err
		}
	case 5:
		{
			switch menu.Userinput {
			case "Y", "y":
				{
					res, err := sqltranslation.PostTask(db, *task)
					if err != nil {
						return "", err
					}
					menu.Userinput = ""
					menu.Menu0 = 0
					menu.Stage = 0
					menu.Allowuserinput = false
					return res, err
				}
			case "E", "e":
				{
					menu.Userinput = ""
					menu.Menu0 = 0
					menu.Stage = 0
					menu.Allowuserinput = false
					return "exiting to main menu", nil
				}
			}
			return "Y for save task, E for exit without saving", nil
		}
	}

	menu.Menu0 = 0
	return "", nil
}

func Update(menu *Switchertype, task *sqltranslation.Tasktable, db *sqlx.DB) (string, error) {
	switch menu.Stage {
	case 1:
		{
			var err error
			task.ID, err = strconv.ParseInt(menu.Userinput, 10, 64)
			if err != nil {
				menu.Stage = 0
				return "", err
			}
			out, err := sqltranslation.GetOverID(db, task)
			if err != nil {
				menu.Stage = 0
				return "", err
			}
			menu.Allowuserinput = false
			menu.Stage++
			return out, nil
		}
	case 2:
		{
			blablabla := `-=UPDATE TASK=-
		0: save and exit to main menu
		1: exit without saving
		2: update message
		3: update status
		4: update deadline
		select option and hit enter:`
			menu.Stage++
			return fmt.Sprintf("task:\n status:%v, deadline:%v, message:%v\n%v\n",
				task.Status, task.Deadline.Format("02.01.06 15.04"), task.Message, blablabla), nil
		}
	case 3:
		{
			switch menu.Userinput {
			case "0":
				{
					res, err := sqltranslation.EditTask(db, *task)
					if err != nil {
						return "", err
					}
					menu.Userinput = ""
					menu.Menu0 = 0
					menu.Allowuserinput = false
					return res, err
				}
			case "1":
				{
					menu.Allowuserinput = false
					menu.Menu0 = 0
					return "exiting to main", nil
				}
			case "2":
				{
					menu.Stage = 4
					return "enter new message:", nil
				}
			case "3":
				{
					menu.Stage = 5
					return "enter new status:", nil
				}
			case "4":
				{
					menu.Stage = 6
					return fmt.Sprintf("enter time of deadline (enter for keep %v)", task.Deadline.Format("15.04")), nil
				}
			case "":
				{
					menu.Allowuserinput = false
					menu.Stage = 2
					return "empty input. Try again", nil
				}
			}
			menu.Allowuserinput = false
			menu.Stage = 2
			return "bad input, try again", nil
		}
	case 4:
		{ // message upd
			if len(menu.Userinput) < 1 {
				return "Empty message entered, please try again or exit", nil
			}
			menu.Stage = 2
			task.Message = menu.Userinput
			menu.Allowuserinput = false
			return "Task message updated", nil
		}
	case 5:
		{ // status upd
			if len(menu.Userinput) < 1 {
				return "Empty status entered, please try again or exit", nil
			}
			menu.Stage = 2
			task.Status = menu.Userinput
			menu.Allowuserinput = false
			return "Task status updated", nil
		}
	case 6:
		{ // time upd
			str, err := UpdateTime(menu, task)
			if err != nil {
				menu.Stage = 2
				menu.Allowuserinput = false
				return "", err
			}
			menu.Stage++
			return fmt.Sprintf("%v\nenter date of deadline (enter for keep %v)", str, task.Deadline.Format("02.01.06")), err
		}
	case 7:
		{ // date upd
			menu.Allowuserinput = false
			menu.Stage = 2
			str, err := UpdateDate(menu, task)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%v\nDeadline updated: %v\n", str, task.Deadline.Format("02.01.06 15:04")), err
		}
	}
	return "", fmt.Errorf("wtf")
}

func Del(menu *Switchertype, task *sqltranslation.Tasktable, db *sqlx.DB) (string, error) {
	switch menu.Stage {
	case 0:
		{
			menu.Stage++
			return "enter task ID for delete:", nil
		}
	case 1:
		{
			var err error
			task.ID, err = strconv.ParseInt(menu.Userinput, 10, 64)
			if err != nil {
				menu.Stage = 0
				return "", err
			}
			out, err := sqltranslation.GetOverID(db, task)
			if err != nil {
				menu.Allowuserinput = false
				menu.Menu0 = 0
				return "", err
			}
			menu.Allowuserinput = false
			menu.Stage++
			return out, nil
		}
	case 2:
		{
			blablabla := "Are you sure you want to delete the task? y for delete, e for exit"
			menu.Stage++
			return fmt.Sprintf("ID:%v, status:%v, deadline:%v, message:%v\n%v\n",
				task.ID, task.Status, task.Deadline.Format("02.01.06 15:04"), task.Message, blablabla), nil
		}
	case 3:
		{
			switch menu.Userinput {
			case "y", "Y":
				{
					menu.Userinput = ""
					menu.Menu0 = 0
					menu.Allowuserinput = false
					return sqltranslation.DelTask(db, *task)
				}
			case "e", "E":
				{
					menu.Userinput = ""
					menu.Menu0 = 0
					menu.Allowuserinput = false
					return "Deletion aborted", nil
				}
			}
			menu.Stage--
			return fmt.Sprintf("%v is bad answer, try again or exit", menu.Userinput), nil
		}
	}

	return "", fmt.Errorf("wtf")
}

func GetWithStatus(menu *Switchertype, db *sqlx.DB) (string, error) {
	switch menu.Stage {
	case 0:
		{
			menu.Stage++
			return "enter task status (empty for 'closed')", nil
		}
	case 1:
		{
			if len(menu.Userinput) < 1 {
				menu.Userinput = "closed"
			}
			out, err := sqltranslation.GetWithStatus(db, menu.Userinput)
			if err != nil {
				menu.Allowuserinput = false
				menu.Menu0 = 0
				return "", err
			}
			menu.Userinput = ""
			menu.Menu0 = 0
			menu.Allowuserinput = false
			return out, nil
		}
	}

	return "wtf", fmt.Errorf("wtf")
}

func GetWithoutStatus(menu *Switchertype, db *sqlx.DB) (string, error) {
	switch menu.Stage {
	case 0:
		{
			menu.Stage++
			return "enter task status (empty for 'closed')", nil
		}
	case 1:
		{
			if len(menu.Userinput) < 1 {
				menu.Userinput = "closed"
			}
			out, err := sqltranslation.GetWithoutStatus(db, menu.Userinput)
			if err != nil {
				menu.Allowuserinput = false
				menu.Menu0 = 0
				return "", err
			}
			menu.Userinput = ""
			menu.Menu0 = 0
			menu.Allowuserinput = false
			return out, nil
		}
	}

	return "wtf", fmt.Errorf("wtf")
}

func Mainmenu(menu *Switchertype) (string, error) {
	menu.Stage = 0
	switch menu.Userinput {
	case "":
		{
			a := (`
	-=MAIN MENU=-
	0: exit
	1: get all tasks 
		11:get tasks with status
		12:get tasks without status
		13:get all deadlined tasks    
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
			menu.Allowuserinput = false
			return "", nil
		}
	case "3":
		{
			menu.Menu0 = 3
			menu.Stage = 1
			return "enter task ID for editing:", nil
		}
	case "4":
		{
			menu.Menu0 = 4
			menu.Allowuserinput = false
			return "", nil
		}
	case "11":
		{
			menu.Menu0 = 11
			menu.Allowuserinput = false
			return "", nil
		}
	case "12":
		{
			menu.Menu0 = 12
			menu.Allowuserinput = false
			return "", nil
		}
	case "13":
		{
			menu.Menu0 = 13
			menu.Allowuserinput = false
			return "", nil
		}
	}
	menu.Allowuserinput = false
	return fmt.Sprintf("Something very strange was entered: %v.\nEnter an existing option.", menu.Userinput), nil
}

func UpdateDate(menu *Switchertype, task *sqltranslation.Tasktable) (string, error) {
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

func UpdateTime(menu *Switchertype, task *sqltranslation.Tasktable) (string, error) {
	if len(menu.Userinput) == 0 {
		return fmt.Sprintf("entered nothing, keeping time=%v", task.Deadline.Format("15.04")), nil
	}
	if len(menu.Userinput) != 5 {
		return "", fmt.Errorf("entered bad sequence: %v, need hh.mm as 15.04", menu.Userinput)
	}
	now, err := time.Parse("02.01.0615.04", task.Deadline.Format("02.01.06")+menu.Userinput)
	if err != nil {
		return "", err
	}
	task.Deadline = now
	return "date updated without any errors", nil
}
