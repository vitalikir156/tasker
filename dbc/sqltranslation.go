package sqltranslation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgress support
)

type Tasktable struct {
	ID       int64     `json:"id"`
	Message  string    `json:"msg"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
	Deadline time.Time `json:"dl"`
}

func Start(config string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", config)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(6)
	db.SetMaxIdleConns(6)
	return db
}

func PostTask(db *sqlx.DB, q Tasktable) (string, error) {
	if len(q.Status) < 1 {
		q.Status = "new"
	}
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	defer func() {
		err2 := tx.Rollback()
		err = errors.Join(err, err2)
	}()
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	qwer := "INSERT INTO tasker.tasker (message, status, deadline) VALUES ($1, $2, $3)"
	out, err := tx.Exec(qwer, q.Message, q.Status, q.Deadline)
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	out2, err := out.RowsAffected()
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	err = tx.Commit()
	return fmt.Sprintf("%v line(s) inserted", out2), err
}

func EditTask(db *sqlx.DB, q Tasktable) (string, error) {
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	defer func() {
		err2 := tx.Rollback()
		err = errors.Join(err, err2)
	}()

	if len(q.Message) > 0 {
		_, err = tx.Exec("update tasker.tasker set message=$1 where id=$2", q.Message, q.ID)
		if err != nil {
			return fmt.Sprintf("%v", err), err
		}
	}
	if len(q.Status) > 0 {
		_, err = tx.Exec("update tasker.tasker set status=$1 where id=$2", q.Status, q.ID)
		if err != nil {
			return fmt.Sprintf("%v", err), err
		}
	}

	if !q.Deadline.IsZero() {
		_, err = tx.Exec("update tasker.tasker set deadline=$1 where id=$2", q.Deadline, q.ID)
		if err != nil {
			return fmt.Sprintf("%v", err), err
		}
	}
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	err = tx.Commit()
	return "update with no errors", err
}

func DelTask(db *sqlx.DB, q Tasktable) (string, error) {
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	defer func() {
		err2 := tx.Rollback()
		err = errors.Join(err, err2)
	}()
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	out, err := tx.Exec("delete from tasker.tasker where id=$1", q.ID)
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	out2, err := out.RowsAffected()
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	return fmt.Sprintf("%v line(s) deleted", out2), err
}

func GetAll(db *sqlx.DB) (string, error) {
	query := `select id,message,status,created,deadline from tasker.tasker`
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var out strings.Builder
	for rows.Next() {
		var task Tasktable
		if err = rows.Scan(&task.ID, &task.Message, &task.Status, &task.Created, &task.Deadline); err != nil {
			fmt.Println(err)
		} else {
			out.WriteString(fmt.Sprintf("ID:%v Status:%v Created:%v Deadline:%v Task:%v \n",
				task.ID, task.Status, task.Created.Format("02.01.06 15:04"), task.Deadline.Format("02.01.06 15:04"), task.Message))
		}
	}
	return out.String(), nil
}

func GetWithStatus(db *sqlx.DB, x string) (string, error) {
	query := `select id,message,status,created,deadline from tasker.tasker where status=$1`
	rows, err := db.QueryContext(context.Background(), query, x)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var out strings.Builder
	for rows.Next() {
		var task Tasktable
		if err = rows.Scan(&task.ID, &task.Message, &task.Status, &task.Created, &task.Deadline); err != nil {
			fmt.Println(err)
		} else {
			out.WriteString(fmt.Sprintf("ID:%v Status:%v Created:%v Deadline:%v Task:%v \n",
				task.ID, task.Status, task.Created.Format("02.01.06 15:04"), task.Deadline.Format("02.01.06 15:04"), task.Message))
		}
	}
	return out.String(), nil
}

func GetWithoutStatus(db *sqlx.DB, s string) (string, error) {
	query := `select id,message,status,created,deadline from tasker.tasker where status<>$1`
	rows, err := db.QueryContext(context.Background(), query, s)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var out strings.Builder
	for rows.Next() {
		var task Tasktable
		if err = rows.Scan(&task.ID, &task.Message, &task.Status, &task.Created, &task.Deadline); err != nil {
			fmt.Println(err)
		} else {
			out.WriteString(fmt.Sprintf("ID:%v Status:%v Created:%v Deadline:%v Task:%v \n",
				task.ID, task.Status, task.Created.Format("02.01.06 15:04"), task.Deadline.Format("02.01.06 15:04"), task.Message))
		}
	}
	return out.String(), nil
}

func GetAllDeadlined(db *sqlx.DB) (string, error) {
	query := `select id,message,status,created,deadline from tasker.tasker where deadline<created`
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var out strings.Builder
	for rows.Next() {
		var task Tasktable
		if err = rows.Scan(&task.ID, &task.Message, &task.Status, &task.Created, &task.Deadline); err != nil {
			fmt.Println(err)
		} else {
			out.WriteString(fmt.Sprintf("ID:%v Status:%v Created:%v Deadline:%v Task:%v \n",
				task.ID, task.Status, task.Created.Format("02.01.06 15:04"), task.Deadline.Format("02.01.06 15:04"), task.Message))
		}
	}
	return out.String(), nil
}

func GetOverID(db *sqlx.DB, q *Tasktable) (string, error) {
	query := `select id,message,status,created,deadline from tasker.tasker where id=$1`
	row := db.QueryRowContext(context.Background(), query, q.ID)

	err := row.Scan(&q.ID, &q.Message, &q.Status, &q.Created, &q.Deadline)
	if err != nil {
		return "", err
	}

	return "get over id: success", nil
}
