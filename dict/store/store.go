package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func InitDbConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&rejectReadOnly=%t&parseTime=true&loc=Local",
		"user",
		"pass",
		"127.0.0.1",
		"3306",
		"ldoce",
		30,
		10,
		10,
		true,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error %s when opening DB\n", err)
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 60)
	return db, nil
}
