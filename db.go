package main

import (
    "log"
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"
)

type PostgresConnection struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func dbConnection() *sql.DB {

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", ConnectionSecrets.Host, ConnectionSecrets.Port, ConnectionSecrets.User, ConnectionSecrets.Password, ConnectionSecrets.DBName)

	/* validate connection */
	db, err := sql.Open("postgres", connectionString)
	if err != nil { log.Panic(err) }

	err = db.Ping()
	if err != nil { log.Panic(err) }

	return db
}

// func AddStudent() (int, error) {

// }

/* privledge enum
*  CREATE TYPE privledge_enum AS ENUM ('STUDENT', 'INSTRUCTOR')
*/

/* initializes db with all the tables */
func DBInit() error {
    db := dbConnection()
    defer db.Close()

    query := `
    CREATE TABLE IF NOT EXISTS users(
        id serial primary key,
        discordid integer not null unique,
        privledge privledge_enum not null
    )
    `
    _, err := db.Exec(query)
    return err
}

/* wipes the db */
func DBPurge() {

}

func DBNewStudent(discordid int) (int, error) {
    db := dbConnection()
    defer db.Close()

    query := `
    INSERT INTO users (discordid, privledge)
    VALUES ($1, $2)
    RETURNING id
    `

    id := -1
    err := db.QueryRow(query, discordid, "STUDENT").Scan(&id)
    return id, err
}

