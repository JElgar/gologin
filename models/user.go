package models

import(
    "database/sql"
)

type UserStore interface {
    GetUser(u *User) (User, *ApiError)
}

type User struct {
    Username    string
    Name        string
    Password    string
}

func (db *DB) GetUser(iu *User) (User, *ApiError){
    // Do a transaction thing to get the user from the database... Need a good way of doing transactions
    var u User
    sqlStmt := `SELECT username, password FROM users WHERE username = $1;`

    if iu == nil {
        return u, &ApiError{nil, "Cannot get NIL user", 400}
    } else if iu.Username == "" {
        return u, &ApiError{nil, "Username Required", 400}
    }


    res := db.QueryRow(sqlStmt, iu.Username)
    switch err := res.Scan(&u.Username, &u.Password); err {
        case sql.ErrNoRows:
            return u, &ApiError{err, "User does not exist in database", 404}
        case nil:
            return u, nil
        default:
            return u, &ApiError{err, "Unknown Error during Query", 400}
    }
}

func (db *DB) CreateUser(u *User) (error) {
    // Add user to database --> again need transation thing
    return nil
}
