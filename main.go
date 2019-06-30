package main

import (
//   "database/sql"
   models "github.com/jelgar/login/models"
   "log"
   "fmt"
   "github.com/gin-gonic/gin"
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
    r := gin.Default()
    r.GET("/ping", ping)
    r.GET("/user", env.getUser)
    r.GET("/createUser", env.createUser)
    r.Run(":8080")
}

func ping(c *gin.Context){
    c.JSON(200, gin.H{
        "world": "Hello",
    })
}

func (e *Env) getUser (c *gin.Context) {

    // This is not how transactions work. Transactions should only return stuff when commited rihgt?
    var u models.User
    c.BindJSON(&u)
    //user, err := e.db.GetUser(&models.User{Username: "john", Password:"123"})
    user, err := e.db.GetUser(&u)
    if err != nil {
        if (err.Code == 404){
            fmt.Println("User doesn't exist")
        }else {
            panic(err)
        }
    }
    fmt.Println(user)
}

func (e *Env) createUser (c *gin.Context){
    // TODO on success return user and enventually JSON web token
    var u models.User
    c.BindJSON(&u)

    user, err := e.db.CreateUser(&u)
    if err != nil && err.Code == 409 {
        fmt.Print("User already exists")
        // TODO Deal with case of collision --> this error code is currently coming out wrong (is 500 should be 409 plz fix 
    } else if err != nil {
        fmt.Println(err.Message)
        panic(err)
    }
    fmt.Println(user)
}
