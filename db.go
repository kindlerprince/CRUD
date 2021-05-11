package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	CREATE = "INSERT INTO UserDetails (name, email , password, address) values ($1, $2, $3, $4) RETURNING name;"
	READ   = "SELECT email,password FROM UserDetails where email=$1 and password=$2;"
	UPDATE = "UPDATE UserDetails SET password=$2 where email=$1 RETURNING email;"
	DELETE = "DELETE FROM UserDetails where email=$1"
)

var (
	createStmt *sql.Stmt
	readStmt   *sql.Stmt
	updateStmt *sql.Stmt
	deleteStmt *sql.Stmt
	dbClient   *sql.DB
)

func dbConnect() error {
	dbAddress := "127.0.0.1"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := "login"
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=5", dbAddress, dbPort, dbUser, dbPassword, dbName)
	dbClient, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("Error in connecting to database %s\n", err.Error())
		return err
	}
	for i := 0; i < 5; i++ {
		if err = dbClient.Ping(); err != nil {
			t := 5
			fmt.Printf("Error in pinging the database : [%s] retrying in %d seconds ...\n", err.Error(), t)
			time.Sleep(time.Duration(t) * time.Second)
		} else {
			fmt.Println("Database connection established")
			break
		}
	}
	createStmt, err = dbClient.Prepare(CREATE)
	if err != nil {
		fmt.Printf("Error in preparing create query %s\n", err.Error())
		return err
	}
	readStmt, err = dbClient.Prepare(READ)
	if err != nil {
		fmt.Printf("Error in preparing read query %s\n", err.Error())
		return err
	}
	updateStmt, err = dbClient.Prepare(UPDATE)
	if err != nil {
		fmt.Printf("Error in preparing update query %s\n", err.Error())
		return err
	}
	deleteStmt, err = dbClient.Prepare(DELETE)
	if err != nil {
		fmt.Printf("Error in preparing delete query %s\n", err.Error())
		return err
	}
	return nil
}
