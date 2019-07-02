package models

import(
    "database/sql"
    "fmt"
    "crypto/rand"
)

type UserStore interface {
    GetUser(u *User) (FullUser, *ApiError)
    CreateUser(u *User) (FullUser, *ApiError)
    Login(u *User) (FullUser, *ApiError)
}

// What types of users do I need?
type User struct {
    Username    string  `json:"username"`
    Password    string  `json:"password"`
}

type FullUser struct {
    Username    string  `json:"username"`
    Name        string  `json:"name"`
    Password    string  `json:"password"`
}

func (db *DB) GetUser(iu *User) (FullUser, *ApiError){
    // Do a transaction thing to get the user from the database... Need a good way of doing transactions
    var u FullUser
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

func (db *DB) CreateUser(u *User) (FullUser, *ApiError) {
    fmt.Println("Entered Create user")
    // TODO Should probablby write function that takes statement and struct and bind
    var fu FullUser
    sqlStmt := `INSERT INTO users (username, password, emailverif, token) VALUES($1,$2, 0, $3);`

    // TODO Verify User is Valid
    _ , err := db.GetUser(u)
    // If the user already exists
    if err == nil {
        // Error code 409 - conflict
       return fu, &ApiError{nil, "The user already exists in the database", 409}
    } else if err.Code != 404 {
        // If there is an error but it is not because the user does not exist
        return fu, &ApiError{err, "Error chechking if user exists", 500}
    }

    // Salt and hash the password (I know this isnt encyption but enc seemed like a nice name)
    encPassword, err := HashSaltPwd([]byte(u.Password))
    if err != nil {
        return fu, &ApiError{err, "Failed to salt and hash password", 500}
    }
    
    // Get a token for email verification
    token, erro := db.getUniqueToken()
    if erro != nil {
        panic(erro)
    }

    res, insertErr := db.Exec(sqlStmt, u.Username, encPassword, token)
    switch insertErr{
        case nil:
            fmt.Println("User inserted")
            fmt.Println(res)
            return fu, nil
        default:
            return fu, &ApiError{err, "Unknown Error during Insertion of User", 400}
    }
}

func (db *DB) Login (u *User) (FullUser, *ApiError) {
    dbUser, err := db.GetUser(u)
    if err != nil {
        // DO THE CHECKS
    }
    isSame, err := ComparePassword(dbUser.Password, []byte(u.Password))
    if err != nil {
        // DO the thing
    }
    if (!isSame) {
        fmt.Println("Incorrect Pssword")
        return dbUser, &ApiError{nil, "Password does not match", 401}
    }

    fmt.Print("Encrypted Passwords Match")
    return dbUser, nil
}

// Generate a random token to assign to a user for password reset or email verification. Generates random token and assures it is unique in database
func (db *DB) getUniqueToken() ([]byte, error){
    unique := false
    var key []byte
    var err error

    for !unique {
        key, err = getToken()
        if err != nil {
            return nil, err
        }
        unique, err = db.tokenIsUnique(key)
        if err != nil {
            return nil, err
        }
    }
    return key, nil
}

func (db *DB) tokenIsUnique (token []byte) (bool, error) {
    unique, err := db.IsUnique(token, "users", "token")
    if err != nil{
        return false, err
    }
    if !unique {
        return false, nil
    }
    return true, nil
}

func getToken() ([]byte, error) {
    token := make([]byte, 20)
    _, err := rand.Read(token)
    if err !=  nil {
        return nil, err
    }
    return token, err
}
