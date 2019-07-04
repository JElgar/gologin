package main

import (
//   "database/sql"
   config "github.com/jelgar/login/config"
   models "github.com/jelgar/login/models"
   "log"
)

// Stuct to store environment variables for application
type Env struct {
    db models.Datastore
}

func main() {

    config.InitConfig()

    db, err := models.InitDB("postgresql://admin:test123@ec2-35-178-198-24.eu-west-2.compute.amazonaws.com/secta?sslmode=disable")
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db: db}

    //Gin stuff
    r := SetupRouter(env)
    r.Run(":8080")
}
