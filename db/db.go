package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// create a db connection
func CreateConn() (*sql.DB, error) {
	var db *sql.DB
	var err error

	host := "127.0.0.1"
	port := "5254"
	user := "postgres"
	password := "brian"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)
	db, err = sql.Open("postgres", psqlInfo)
	defer db.Close()
	if err != nil {
		log.Fatalf("Couldnt connect due to : %s", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed ping: %s", err)
		return nil, err
	}

	err = CreateDB()
	if err != nil {
		log.Fatalf("Failed to create db: %s", err)
		return nil, err
	}
	_, err = db.Exec("USE employee")
	if err != nil {
		log.Fatalf("Failed to change db: %s", err)
		return nil, err
	}
	log.Printf("\nSuccessfully connected to database!\n")

	return db, nil

}

// create a db transaction
func transaction() (*sql.Tx, error) {
	db, err := CreateConn()
	if err != nil {
		log.Fatalf("Couldnt connect due to : %s", err)
		return nil, err
	}

	// background context is fine
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Couldnt create a transaction due to : %s", err)
		return nil, err
	}

	return tx, nil
}

// CreateDB creates a db used by the app
func CreateDB() error {
	tx, err := transaction()
	error_handle(err)
	_, err = tx.Exec(`SELECT 'CREATE DATABASE employees'
					  WHERE NOT EXISTS 
					  (SELECT FROM pg_database WHERE datname = 'employee')`)
	error_handle(err)
	log.Println("Created db successfully")
	return nil
}

// error handling
func error_handle(err error) error {
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}
