package models

import(
    "database/sql"
    "fmt"
    errors "github.com/jelgar/login/errors"
)

type UserStore interface {
    GetUser(u *User) (User, *errors.ApiError)
    CreateUser(u *User) (User, *errors.ApiError)
    Login(u *User) (User, *errors.ApiError)
    VerfUserEmail(token string) *errors.ApiError
    PasswordRest(token string) *errors.ApiError
}

// What types of users do I need?
type User struct {
    Username    string  `json:"username"`
    Password    string  `json:"password"`
    Email       string  `json:"email"`
    Name        string  `json:"name"`
    EmailToken  string  `json:"emailverif"`
}

func (db *DB) GetUser(iu *User) (User, *errors.ApiError){
    // Do a transaction thing to get the user from the database... Need a good way of doing transactions
    var u User
    sqlStmt := `SELECT username, password FROM users WHERE username = $1;`

    if iu == nil {
        return u, &errors.ApiError{nil, "Cannot get NIL user", 400}
    } else if iu.Username == "" {
        return u, &errors.ApiError{nil, "Username Required", 400}
    }


    res := db.QueryRow(sqlStmt, iu.Username)
    switch err := res.Scan(&u.Username, &u.Password); err {
        case sql.ErrNoRows:
            return u, &errors.ApiError{err, "User does not exist in database", 404}
        case nil:
            return u, nil
        default:
            return u, &errors.ApiError{err, "Unknown Error during Query", 500}
    }
}

func (db *DB) CreateUser(u *User) (User, *errors.ApiError) {
    fmt.Println("Entered Create user")
    // TODO Should probablby write function that takes statement and struct and bind
    // TODO Find and replace and swap all instances of fu with empty user obj
    var fu User
    sqlStmt := `INSERT INTO users (username, email, password, emailverif, token) VALUES($1,$2, $3, '0', $4);`

    // TODO Verify User is Valid
    _ , err := db.GetUser(u)
    // If the user already exists
    if err == nil {
        // Error code 409 - conflict
       return fu, &errors.ApiError{nil, "The user already exists in the database", 409}
    } else if err.Code != 404 {
        // If there is an error but it is not because the user does not exist
        return fu, &errors.ApiError{err, "Error chechking if user exists", 500}
    }

    // Salt and hash the password (I know this isnt encyption but enc seemed like a nice name)
    encPassword, err := HashSaltPwd([]byte(u.Password))
    if err != nil {
        return fu, &errors.ApiError{err, "Failed to salt and hash password", 500}
    }
    
    // TODO Change this into a transaction so i add the token after so i can a seperate token function maybe? Possily not?

    // Get a token for email verification
    token, erro := db.getUniqueToken()
    if erro != nil {
        return fu, &errors.ApiError{nil, "There was an error gettign a unique token for email verif", 401}
    }

    fmt.Println("Got the unique token")
    // TODO make sure only insert a max number of token items (so it doesnt go above 50)
    res, insertErr := db.Exec(sqlStmt, u.Username, u.Email, encPassword, token)
    switch insertErr{
        case nil:
            fmt.Println("User inserted")
            fmt.Println(res)
            u.EmailToken = token
            db.SendVerfEmail(u)
            return fu, nil
        default:
            return fu, &errors.ApiError{err, "Unknown Error during Insertion of User", 400}
    }
}


func (db *DB) Login (u *User) (User, *errors.ApiError) {
    dbUser, err := db.GetUser(u)
    if err != nil {
        switch err.Code {
            case 400:
                // This should not be 400 this is a definately the wrong number
                // But then it kind of could be an internal server error cause front end should never allow the input of a empty ting so idk but then no this isnt so 400 is so wrong
                return dbUser, &errors.ApiError{nil, "Invalid input of user credentials when checking if user existed during login", 400}
            case 404:
                return dbUser, &errors.ApiError{nil, "User that is trying to login does not exist in database", 404}
            default:
                return dbUser, &errors.ApiError{nil, "Uknown error when trying to check if loging in user exists", 500}
        }
    }
    isSame, err := ComparePassword(dbUser.Password, []byte(u.Password))
    if err != nil {
        return dbUser, &errors.ApiError{nil, "Error comparing users", 401}
    }
    if (!isSame) {
        fmt.Println("Incorrect Pssword")
        return dbUser, &errors.ApiError{nil, "Password does not match", 401}
    }

    fmt.Print("Encrypted Passwords Match")
    return dbUser, nil
}


// TODO Ok im writting this comment to try and work out what im supposed to be doing. So... I have a routers section that calls this section to do stuff. Do i not want to split up these functions into smaller fucntions causes like atm some of them are huge/repetative. But then i feel like i need another layer or something like so i can simply send sql commands/transaction to the db and have seperate thing for idk ... Im lost ??? Maybe i will jsut keep doing what im doing but if it get repetative im gonna need to change plans
// Oh the reason i had this moment now was because either i could write settoken or i could do a whole passwrod reset function.... surely it better to split up right?

// TODO also add these extra function to the UserDatastroe interface
func (db *DB) SetToken (u *User) (*errors.ApiError) {
    return nil
}

func (db *DB) PasswordRest (string) (*errors.ApiError) {
    return nil
}
