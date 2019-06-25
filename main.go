package main

import (
//   "database/sql"
   models "github.com/jelgar/login/models"
   "log"
   "fmt"
)

// Put this is models later
type Env struct {
    db models.Datastore
}

func main() {
    db, err := models.InitDB("postgresql://admin:test123@ec2-35-178-198-24.eu-west-2.compute.amazonaws.com/secta?sslmode=disable")
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db: db}

    fmt.Println(env)

    //Gin stuff
}
