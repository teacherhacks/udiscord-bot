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


/* initializes db with all the tables */
func InitDB() error {
    db := dbConnection()
    defer db.Close()

    query := `
    CREATE TABLE IF NOT EXISTS users(
        id serial primary key,
        discordid varchar(64) not null unique,
        privledge privledge_enum not null
    )
    `
    _, err := db.Exec(query)
    return err
}

/* wipes the db */
func PurgeDB() {

}

