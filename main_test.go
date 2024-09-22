package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	sqltranslation "github.com/vitalikir156/tasker/dbc"
	ui "github.com/vitalikir156/tasker/userinterface"
)

const (
	userinput1 = "test message!"
	userinput2 = "teststatus"
	userinput3 = "17.30"
	userinput4 = "17.03.98"
)

func TestDBCPostGood(t *testing.T) {
	// по сути все кейсы DBC покрыты в интеграционном тестировании, тут только проба пера
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockdb.Close()
	db := sqlx.NewDb(mockdb, "sqlmock")
	time, _ := time.Parse("02.01.0615.04", "17.03.9817.30")
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tasker.tasker").WithArgs("TestMsg1", "new", time).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	task := sqltranslation.Tasktable{Message: "TestMsg1", Status: "new", Deadline: time}
	out, err := sqltranslation.PostTask(db, task)
	require.Equal(t, out, "1 line(s) inserted")
	require.NoError(t, err)
}

func TestUIPostGood(t *testing.T) { // проверяет корректность набивания структуры данными (Post)
	mockdb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockdb.Close()
	db := sqlx.NewDb(mockdb, "sqlmock")
	userinput1 := "test message!"
	userinput2 := "teststatus"
	userinput3 := "17.30"
	userinput4 := "17.03.98"
	var selector ui.Switchertype
	var task sqltranslation.Tasktable
	out, err := ui.Post(&selector, &task, db)
	require.NoError(t, err)
	require.Equal(t, out, "enter task message:")
	require.Equal(t, selector.Stage, 1)
	selector.Userinput = userinput1
	out, err = ui.Post(&selector, &task, db)
	require.NoError(t, err)
	require.Equal(t, out, "enter status(empty for <new>):")
	require.Equal(t, selector.Stage, 2)
	selector.Userinput = userinput2
	_, err = ui.Post(&selector, &task, db)
	require.NoError(t, err)
	require.Equal(t, selector.Stage, 3)
	selector.Userinput = userinput3
	_, err = ui.Post(&selector, &task, db)
	require.NoError(t, err)
	require.Equal(t, selector.Stage, 4)
	selector.Userinput = userinput4
	_, err = ui.Post(&selector, &task, db)
	require.NoError(t, err)
	require.Equal(t, selector.Stage, 5)

	require.Equal(t, task.Message, userinput1)
	require.Equal(t, task.Status, userinput2)
	require.Equal(t, task.Deadline.Format("02.01.06 15.04"), fmt.Sprintf("%v %v", userinput4, userinput3))
}

func TestUIUpdateGood(t *testing.T) { // проверяет корректность набивания структуры данными (Update)
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockdb.Close()
	var rowsql sqltranslation.Tasktable
	testid := 1
	rows := sqlmock.NewRows([]string{"id", "message", "status", "created", "deadline"}).
		AddRow(rowsql.ID, "rowsql.Message", rowsql.Status, rowsql.Created, rowsql.Deadline)
	mock.ExpectQuery("select id,message,status,created,deadline from tasker.tasker where").
		WithArgs(testid).WillReturnRows(rows)

	db := sqlx.NewDb(mockdb, "sqlmock")
	var selector ui.Switchertype
	var task sqltranslation.Tasktable
	userinput0 := "1"

	selector.Stage = 1
	selector.Userinput = userinput0
	_, err = ui.Update(&selector, &task, db)
	require.NoError(t, err)
	selector.Stage = 3
	selector.Userinput = "2"
	_, err = ui.Update(&selector, &task, db)
	require.NoError(t, err)
	require.Equal(t, selector.Stage, 4)
	selector.Userinput = userinput1
	_, err = ui.Update(&selector, &task, db)
	require.NoError(t, err)
	selector.Stage = 3
	selector.Userinput = "3"
	_, err = ui.Update(&selector, &task, db)
	require.NoError(t, err)
	require.Equal(t, selector.Stage, 5)
	selector.Userinput = userinput2
	_, err = ui.Update(&selector, &task, db)
	require.NoError(t, err)
	selector.Stage = 3
	selector.Userinput = "4"
	_, err = ui.Update(&selector, &task, db)
	require.NoError(t, err)
	selector.Userinput = userinput3
	_, err = ui.Update(&selector, &task, db)
	require.NoError(t, err)
	selector.Userinput = userinput4
	_, err = ui.Update(&selector, &task, db)
	require.NoError(t, err)
	require.Equal(t, task.Message, userinput1)
	require.Equal(t, task.Status, userinput2)
	require.Equal(t, task.Deadline.Format("02.01.06 15.04"), fmt.Sprintf("%v %v", userinput4, userinput3))
}

func TestMainGood1(t *testing.T) { // getters, post
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockdb.Close()
	var rowsql sqltranslation.Tasktable
	// testid:=1
	status := "new"
	time, err := time.Parse("02.01.0615.04", fmt.Sprint(userinput4, userinput3))
	require.NoError(t, err)
	rows := sqlmock.NewRows([]string{"id", "message", "status", "created", "deadline"}).
		AddRow(rowsql.ID, "rowsql.Message", rowsql.Status, rowsql.Created, rowsql.Deadline)
	rows2 := sqlmock.NewRows([]string{"id", "message", "status", "created", "deadline"}).
		AddRow(rowsql.ID, "rowsql.Message", "new", rowsql.Created, rowsql.Deadline)
	rows3 := sqlmock.NewRows([]string{"id", "message", "status", "created", "deadline"}).
		AddRow(rowsql.ID, "rowsql.Message", "notnew", rowsql.Created, rowsql.Deadline)
	rows4 := sqlmock.NewRows([]string{"id", "message", "status", "created", "deadline"}).
		AddRow(rowsql.ID, "deadlined", rowsql.Status, rowsql.Created, rowsql.Deadline)

	mock.ExpectQuery("select id,message,status,created,deadline from tasker.tasker").WillReturnRows(rows)
	mock.ExpectQuery("select id,message,status,created,deadline from tasker.tasker where ").
		WithArgs(status).WillReturnRows(rows2)
	mock.ExpectQuery("select id,message,status,created,deadline from tasker.tasker where ").
		WithArgs(status).WillReturnRows(rows3)
	mock.ExpectQuery("select id,message,status,created,deadline from tasker.tasker where ").
		WillReturnRows(rows4)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tasker.tasker").WithArgs(userinput1, userinput2, time).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	db := sqlx.NewDb(mockdb, "sqlmock")
	var task sqltranslation.Tasktable
	var switcher ui.Switchertype
	out, err := menucore(&switcher, &task, db)
	require.NoError(t, err)
	out1 := `
	-=MAIN MENU=-
	0: exit
	1: get all tasks 
		11:get tasks with status
		12:get tasks without status
		13:get all deadlined tasks    
	2: create task
	3: edit task
	4: delete task
	select option and hit enter:`
	require.Equal(t, out1, out)
	switcher.Userinput = "1"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	out, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	require.Equal(t, "ID:0 Status: Created:01.01.01 00:00 Deadline:01.01.01 00:00 Task:rowsql.Message \n",
		out) // getall result
	switcher.Userinput = "11"
	menucore(&switcher, &task, db)
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = status
	out, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	require.Equal(t, "ID:0 Status:new Created:01.01.01 00:00 Deadline:01.01.01 00:00 Task:rowsql.Message \n",
		out) // getwithstatus
	switcher.Userinput = "12"
	menucore(&switcher, &task, db)
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = status
	out, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	require.Equal(t, "ID:0 Status:notnew Created:01.01.01 00:00 Deadline:01.01.01 00:00 Task:rowsql.Message \n",
		out) // getwithoutstatus
	switcher.Userinput = "13"
	menucore(&switcher, &task, db)
	out, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	require.Equal(t, "ID:0 Status: Created:01.01.01 00:00 Deadline:01.01.01 00:00 Task:deadlined \n", out) // getdeadlined

	switcher.Userinput = "2"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	out, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	require.Equal(t, "enter task message:", out)
	switcher.Userinput = userinput1
	out, err = menucore(&switcher, &task, db) // post message
	require.NoError(t, err)
	require.Equal(t, "enter status(empty for <new>):", out)
	switcher.Userinput = userinput2
	_, err = menucore(&switcher, &task, db) // post status
	require.NoError(t, err)
	switcher.Userinput = userinput3
	_, err = menucore(&switcher, &task, db) // time
	require.NoError(t, err)
	switcher.Userinput = userinput4
	_, err = menucore(&switcher, &task, db) // date
	require.NoError(t, err)
	switcher.Userinput = "y"
	out, err = menucore(&switcher, &task, db) // Y for save to DB
	require.NoError(t, err)
	require.Equal(t, "1 line(s) inserted", out) // Ufff... Its working
}

func TestMainGood2(t *testing.T) { // update and delete
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockdb.Close()
	var rowsql sqltranslation.Tasktable
	time, err := time.Parse("02.01.0615.04", fmt.Sprint(userinput4, userinput3))
	require.NoError(t, err)

	rows5 := sqlmock.NewRows([]string{"id", "message", "status", "created", "deadline"}).
		AddRow(rowsql.ID, "not edited", rowsql.Status, rowsql.Created, rowsql.Deadline)
	rows6 := sqlmock.NewRows([]string{"id", "message", "status", "created", "deadline"}).
		AddRow(rowsql.ID, "not deleted", rowsql.Status, rowsql.Created, rowsql.Deadline)

	mock.ExpectQuery("select id,message,status,created,deadline from tasker.tasker where ").
		WithArgs(0).WillReturnRows(rows5)
	mock.ExpectBegin()
	mock.ExpectExec("update tasker.tasker set").WithArgs(userinput1, 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update tasker.tasker set").WithArgs(userinput2, 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update tasker.tasker set").WithArgs(time, 0).WillReturnResult(sqlmock.
		NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("select id,message,status,created,deadline from tasker.tasker where ").
		WithArgs(0).WillReturnRows(rows6)
	mock.ExpectBegin()
	mock.ExpectExec("delete from tasker.tasker where").WithArgs(0).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	db := sqlx.NewDb(mockdb, "sqlmock")

	var task sqltranslation.Tasktable
	var switcher ui.Switchertype

	require.Equal(t, 0, switcher.Menu0)
	switcher.Userinput = "3"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = "0"
	menucore(&switcher, &task, db)
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = "2"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = userinput1
	_, err = menucore(&switcher, &task, db) // msg
	require.NoError(t, err)
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = "3"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = userinput2
	_, err = menucore(&switcher, &task, db) // status
	require.NoError(t, err)
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = "4"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = userinput3
	_, err = menucore(&switcher, &task, db) // time
	require.NoError(t, err)
	switcher.Userinput = userinput4
	_, err = menucore(&switcher, &task, db) // date
	require.NoError(t, err)
	out, err := menucore(&switcher, &task, db)
	require.NoError(t, err)
	updut := `task:
 status:teststatus, deadline:17.03.98 17.30, message:test message!
-=UPDATE TASK=-
		0: save and exit to main menu
		1: exit without saving
		2: update message
		3: update status
		4: update deadline
		select option and hit enter:
`
	require.Equal(t, updut, out)
	switcher.Userinput = "0"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	require.Equal(t, 0, switcher.Menu0)
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = "4"
	menucore(&switcher, &task, db)
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	switcher.Userinput = "0"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	out, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
	del := `ID:0, status:, deadline:01.01.01 00:00, message:not deleted
Are you sure you want to delete the task? y for delete, e for exit
`
	require.Equal(t, del, out)
	switcher.Userinput = "y"
	_, err = menucore(&switcher, &task, db)
	require.NoError(t, err)
}
