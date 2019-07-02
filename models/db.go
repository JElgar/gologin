package models

import (
    "database/sql"
    _ "github.com/lib/pq"
)

// This interface will contain all the interfaces of the required operations of the Datastore. For example a user Datastore interface will be created and used to define all methods for user based operations.
type Datastore interface {
    UserStore
    Begin() (*Tx, error)
}

// This is a DB sturct that will inherit the interface above and thus MUST implement all the required methods // Is this a good way of doing this??
type DB struct {
    *sql.DB
}

// General transaction 
type Tx struct {
    *sql.Tx
}

func InitDB(DbUrl string) (*DB, error) {
    db, err := sql.Open("postgres", DbUrl)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return &DB{db}, nil
}

func (db *DB) Begin() (*Tx, error) {
    tx, err := db.DB.Begin()
    if err != nil {
        return nil, err 
    }
    return &Tx{tx}, nil
}

// May want to move this into a seperate db package/file at somepoint ?
// Given a table, column and some data, checks if that data already exists in the column
func (db *DB) IsUnique(data interface{}, table string, column string) (bool, error){
    var count int
    sqlStmt := `SELECT COUNT($1) FROM $2 WHERE username = $3;`
    row := db.QueryRow(sqlStmt, column, table, data)
    if err := row.Scan(&count); err != nil {
        return false, err
    }
    if count != 0 {
        return false, nil
    }
    return true, nil
}
