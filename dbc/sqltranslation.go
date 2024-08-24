package sqltranslation

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgress support
)

type tasker struct {
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
func PostTask(db *sqlx.DB, q tasker) (string, error) {

	fmt.Println(q)
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	defer tx.Rollback()
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	out, err := tx.Exec("INSERT INTO tasker.tasker (message, status) VALUES ($1, $2)", q.Message, "new")
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
