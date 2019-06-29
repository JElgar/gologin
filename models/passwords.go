package models

import (
    "golang.org/x/crypto/bcrypt"
)

func HashSaltPwd(pwd []byte) (string, *ApiError) {
    // TODO look into what cost i should use (needs to be > 4)
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        return "", &ApiError{err, "Error hashing and salting password", 500}
    }

    return string(hash), nil
}

func ComparePassword(hashedPwd string, plainPwd []byte) (bool, *ApiError) {
    byteHash := []byte(hashedPwd)

    err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
    if err != nil {
        return false, &ApiError{err, "Error Comapring Passwords", 500}
    }
    return true, nil

}
