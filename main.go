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

    db, err := models.InitDB(config.DbString)
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db: db}

    //Gin stuff
    r := SetupRouter(env)
    r.Run(":8080")
}
