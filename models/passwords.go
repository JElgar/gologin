package models

import (
    "golang.org/x/crypto/bcrypt"
    errors "github.com/jelgar/login/errors"
)

func HashSaltPwd(pwd []byte) (string, *errors.ApiError) {
    // TODO look into what cost i should use (needs to be > 4)
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        return "", &errors.ApiError{err, "Error hashing and salting password", 500}
    }

    return string(hash), nil
}

func ComparePassword(hashedPwd string, plainPwd []byte) (bool, *errors.ApiError) {
    byteHash := []byte(hashedPwd)

    err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
    if err != nil {
        return false, &errors.ApiError{err, "Error Comapring Passwords", 500}
    }
    return true, nil

}
