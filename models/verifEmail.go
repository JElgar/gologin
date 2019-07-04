package models

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    errors "github.com/jelgar/login/errors"
    mail "github.com/jelgar/login/email"
    config "github.com/jelgar/login/config"
)

func (db *DB) SendVerfEmail(u *User) {
    // TODO dont forget to update domain
    url := config.Domain + "/confirmEmail?token=" + u.EmailToken
    mail.Send(u.Username, u.Email, url, "email/email.html")
}

// Given a token compate against all users in database and set verifemial true if token matches
func (db *DB) VerfUserEmail(token string) *errors.ApiError {
    sqlStmt := `UPDATE users SET emailverif='1', token='' WHERE token=$1;`
    res, err := db.Exec(sqlStmt, token)
    if (err != nil) {
        return &errors.ApiError{err, "Could not update email verif", 500}
    }
    ra, err := res.RowsAffected()
    if err != nil {
        return &errors.ApiError{err, "Error chekcing number of rows affected during user verify", 500}
    }
    if (ra == 1){
        fmt.Println("Email Verified")
    } else {
        fmt.Println("No address with this token")
        return &errors.ApiError{nil, "Token not associtated with any account", 404}
        // TODO propograte this error and eventually send with gin
    }

    // TODO if email verified redirect to sign in, if not redirect to 404 
    // What am I saying this is definately a front end thing but porbs need to send jobs a gooden or something so it know it can reidrect?
    // Actaully no need a way to say no cells were changed so need to maybe count user or osomething idk, does SQl return number of cells changed?
    return nil
}

// Generate a random token to assign to a user for password reset or email verification. Generates random token and assures it is unique in database
func (db *DB) getUniqueToken() (string, error){
    unique := false
    var key string
    var err error

    for !unique {
        key, err = getToken(20)
        if err != nil {
            return "", err
        }
        unique, err = db.tokenIsUnique(key)
        if err != nil {
            return "", err
        }
    }
    return key, nil
}

func (db *DB) tokenIsUnique (token string) (bool, error) {
    fmt.Println("Checking if token is unique")
    unique, err := db.IsUnique(token, "users", "token")
    if err != nil{
        return false, err
    }
    if !unique {
        return false, nil
    }
    return true, nil
}

func getToken(length int) (string, error) {
    token := make([]byte, length)
    _, err := rand.Read(token)
    if err !=  nil {
        return "", err
    }

    return base64.URLEncoding.EncodeToString(token), nil
}
