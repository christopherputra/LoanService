package database

import (
	"database/sql"
	"log"
)

func (pg *PostgresServer) StartTransaction() (*sql.Tx, error) {
	tx, err := pg.Db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction: ", err)
		return nil, err
	}
	return tx, nil
}

func (pg *PostgresServer) CommitTransaction(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		log.Fatal("Error committing transaction: ", err)
		return err
	}
	return nil
}
