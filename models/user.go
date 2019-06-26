package models

import(
    "errors"
)

type UserStore interface {
}

type User struct {
    Username    string
    Name        string
    Password    string
}

func (tx *Tx) GetUser(u *User) (*User, error){
    // Do a transaction thing to get the user from the database... Need a good way of doing transactions
    if u == nil {
        return nil, errors.New("user required")
    } else if u.Username == "" {
        return nil, errors.New("Username required")
    }


    _, err := tx.Exec(`INSERT INTO user(username, password) VALUES(?,?)`)
    if err != nil {

    }
    return nil, nil
}

func (tx *Tx) CreateUser(u *User) (error) {
    // Add user to database --> again need transation thing
    return nil
}
