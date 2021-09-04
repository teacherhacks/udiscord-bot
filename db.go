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

type Assignment struct {
    ID      int
    Name    string
    Due     int64
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

/* initializes db with all the tables */
func DBInit() error {
    db := dbConnection()
    defer db.Close()

    query := `
    CREATE TYPE privledge_enum AS ENUM ('STUDENT', 'INSTRUCTOR');
    CREATE TABLE users(
        id serial primary key,
        guildid varchar(64) not null,
        discordid varchar(64) not null,
        privledge privledge_enum not null
    );
    CREATE TABLE assignments(
        id serial primary key,
        guildid varchar(64) not null,
        name varchar(256) not null,
        due bigint not null
    );
    `
    _, err := db.Exec(query)
    return err
}

/* wipes the db */
func DBPurge() {

}

func DBNewStudent(guildID string, discordID string) (int, error) {
    db := dbConnection()
    defer db.Close()

    query := `
    INSERT INTO users (guildid, discordid, privledge)
    VALUES ($1, $2)
    RETURNING id
    `

    id := -1
    err := db.QueryRow(query, guildID, discordID, "STUDENT").Scan(&id)
    return id, err
}

func DBNewAssignment(guildID string, name string, due int64) (int, error) {
    db := dbConnection()
    defer db.Close()

    query := `
    INSERT INTO assignments (guildid, name, due)
    VALUES ($1, $2)
    RETURNING id
    `

    id := -1
    err := db.QueryRow(query, guildID, name, due).Scan(&id)
    return id, err
}

/* get all assignments from a guild */
// func DBGetAssignment(guildID int) ([]*Assignment, error) {
func DBGetAssignment(guildID string) {
    db := dbConnection()
    defer db.Close()

    query := `
    SELECT * FROM assignments WHERE guildid == $1
    `

    _, _ = db.Query(query, guildID)
}

