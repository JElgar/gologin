package models

type UserStore interface {
    GetUser() (User, error)
    CreateUser(User) (User, error)
}

type User struct {
    Username    string
    Name        string
    Password    string
}

func (db *DB) GetUser() (User, error){
    // Do a transaction thing to get the user from the database... Need a good way of doing transactions
    return User{}, nil
}

func (db *DB) CreateUser(User) (User, error) {
    // Add user to database --> again need transation thing
    return User{}, nil
}
