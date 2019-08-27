package models

import (
    "golang.org/x/crypto/bcrypt"
    errors "github.com/jelgar/login/errors"
    "fmt"
)

func HashSaltPwd(pwd []byte) (string, *errors.ApiError) {
    // TODO look into what cost i should use (needs to be > 4)
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        return "", &errors.ApiError{err, "Error hashing and salting password", 500}
    }

    return string(hash), nil
}

// The differences between this fucntion and the internal bcrypt password is dumb
func ComparePassword(hashedPwd string, plainPwd []byte) (bool, *errors.ApiError) {
    byteHash := []byte(hashedPwd)

    err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
    if err == bcrypt.ErrMismatchedHashAndPassword {
        fmt.Println("Passwords do not match")
        return false, nil
    } else if err != nil {
        fmt.Println(err)
        panic(err)
        return false, &errors.ApiError{err, "Error Comapring Passwords", 500}
    }
    return true, nil
}
