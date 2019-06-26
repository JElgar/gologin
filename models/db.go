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
