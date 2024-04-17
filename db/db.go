package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// create a db connection
func createConn() (*sql.DB, error) {
	var db *sql.DB
	var err error

	host := "127.0.0.1"
	port := "5254"
	user := "postgres"
	password := "brian"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Couldnt connect due to : %s", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed ping: %s", err)
		return nil, err
	}

	log.Printf("\nSuccessfully connected to database!\n")
	return db, nil

}

// create a db transaction
func transaction() (*sql.Tx, error) {
	db, err := createConn()
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
	if err != nil {
		log.Fatalf("Couldnt create a transaction due to : %s", err)
		return err
	}

	_, err = tx.Exec(`SELECT 'CREATE DATABASE employees'
					  WHERE NOT EXISTS 
					  (SELECT FROM pg_database WHERE datname = 'employee')`)
	if err != nil {
		log.Fatalf("Couldnt create db due to : %s", err)
		return err
	}

	log.Println("Created db successfully")
	return nil
}

// create