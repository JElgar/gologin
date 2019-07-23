package models

import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
)

// This interface will contain all the interfaces of the required operations of the Datastore. For example a User Datastore interface will be created and used to define all methods for user based operations.
type Datastore interface {
    UserStore
    Begin() (*Tx, error)
}

// This is a DB sturct that will inherit the interface above and thus MUST implement all the required methods // Is this a good way of doing this??
type DB struct {
    *sql.DB
}

// General transaction 
// TODO havent really bothered to sort this out yet, so probably look into this
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
// If you dont hate everything dont do this...
// TODO check if this actually works
func (db *DB) IsUnique(data interface{}, table string, column string) (bool, error){
    var count int
    var sqlStmt string
    var row *sql.Row
    sqlStmt = "SELECT COUNT(*) FROM " + table + " WHERE " + column + " = $1;"

    switch data.(type) {
        case string:
            fmt.Println("got a string")
            var d string
            d = data.(string)
            fmt.Println(d)
            row = db.QueryRow(sqlStmt, d)
            fmt.Println("got a row")
        case int:
            var d int
            d = data.(int)
            row = db.QueryRow(sqlStmt, d)
    }

    if err := row.Scan(&count); err != nil {
        return false, err
    }
    if count != 0 {
        return false, nil
    }
    fmt.Println("The token is unique happy days")
    return true, nil
}
