package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	query :=
		`
				INSERT INTO dicts (word, content)
				VALUES (?, ?)
				ON DUPLICATE KEY UPDATE word = VALUES(word), content = VALUES(content)
	`

	// Open a connection to the database
	db, err := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/ldoce")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	result, err := db.Exec(query, "kaixin-test", "hello_world")
	if err != nil {
		log.Fatal(err)
	}

	id, _ := result.LastInsertId()
	affected, _ := result.RowsAffected()
	log.Println("Insert successful", id, affected)
}
