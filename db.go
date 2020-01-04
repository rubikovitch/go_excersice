package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var host = getenv("PSQL_HOST", "localhost")
var port = getenv("PSQL_PORT", "5432")
var user = getenv("PSQL_USER", "postgres")
var password = getenv("PSQL_PWD", "docker")
var dbname = getenv("PSQL_DB_NAME", "postgres")

func saveTimeToDB(time string) {
	db, error := openDb()
	if error != nil {
		panic(error)
	}
	defer db.Close()

	createStatement := `
	CREATE TABLE IF NOT EXISTS time (id SERIAL PRIMARY KEY, time VARCHAR NOT NULL)
	`
	_, error = db.Exec(createStatement)
	if error != nil {
		panic(error)
	}

	insertStatement := `
		INSERT INTO time (time)
		VALUES ($1)
		RETURNING id
	`
	id := 0
	error = db.QueryRow(insertStatement, time).Scan(&id)
	if error != nil {
		panic(error)
	}
	fmt.Println("New time record ID is: ", id)
}

func fetchRecentTime() {
	db, error := openDb()
	if error != nil {
		panic(error)
	}
	defer db.Close()

	latestRowStatement := `
		SELECT id, time 
		FROM time 
		ORDER BY id DESC
		LIMIT 1
	`
	var id int
	var time string
	row := db.QueryRow(latestRowStatement)
	switch err := row.Scan(&id, &time); err {
	case sql.ErrNoRows:
		fmt.Println("No rows are inserted yet")
	case nil:
		fmt.Println("Recent time with ID are: ", time, id)

	default:
		panic(err)
	}
}

func fetchAllTime() {
	db, error := openDb()
	if error != nil {
		panic(error)
	}
	defer db.Close()

	fetchAllStatement := `
		SELECT *
		FROM time
		ORDER BY id DESC
	`
	rows, err := db.Query(fetchAllStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	log.Println("Fetching all stored time")
	for rows.Next() {
		var id int
		var time string
		if err := rows.Scan(&id, &time); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID %v value %v\n", id, time)
	}
}

func openDb() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return sql.Open("postgres", psqlInfo)
}
