package email

import (
    "crypto/rand"
   models "github.com/jelgar/login/models"
)

// Generate a random token to assign to a user for password reset or email verification. Generates random token and assures it is unique in database
func GetUniqueToken(db *models.DB) ([]byte, error){
    unique := false
    var key []byte
    var err error

    for !unique {
        key, err = getToken()
        if err != nil {
            return nil, err
        }
        unique, err = isUnique(key, db)
        if err != nil {
            return nil, err
        }
    }
    return key, nil
}

func isUnique (token []byte, db *models.DB) (bool, error) {
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
