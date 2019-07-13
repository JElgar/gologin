package main

import (
   models "github.com/jelgar/login/models"
   email "github.com/jelgar/login/email"
   "fmt"
   "github.com/gin-gonic/gin"
)

func SetupRouter(env *Env) *gin.Engine {
    r := gin.Default()
    r.GET("/ping", ping)
    r.GET("/user", env.getUser)
    r.POST("/createUser", env.createUser)
    r.POST("/login", env.login)
    r.POST("/sendMail", env.sendMail)
    r.GET("/confirmEmail", env.confirmEmail)

    return r
}


func ping(c *gin.Context){
    c.JSON(200, gin.H{
        "world": "Hello",
    })
}

// Get user accepts a JSON object contains the username of the user it wishes to find
// Get this working for email
func (e *Env) getUser (c *gin.Context) {
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
        // Actauly may be wokring need to check
    } else if err != nil {
        fmt.Println(err.Message)
        panic(err)
    }
    fmt.Println(user)
}

func (e *Env) login (c *gin.Context) {
    var u models.User
    c.BindJSON(&u)

    user, err := e.db.Login(&u)
    if err != nil {
        // TODO return the correct stuff here
        // Ie return actaul json with gin dont just print some random stuff out
        switch err.Code {
            // 401 -> Incorrect password
            case 401:
                fmt.Println("Incorrect Password")
            case 404:
                fmt.Println("User does not exist")
            case 500:
                fmt.Println("Uknown error so so sorry")
            default:
                fmt.Println("retunr a 500 -> Unknown error")
        }
    }
    fmt.Println(user)
    // If user exists return a JWT being like yup and err nill
    // Otherwise return no JWT and be like that this was the error -> eg no user
}

func (e *Env) sendMail (c *gin.Context) {
    // This is a test handler to send emails to a user
    err := email.Send("James", "jamezy850@gmail.com", "jameselgar.com", "email/email.html")
    if err != nil {
        panic(err) 
    }

}

func (e *Env) confirmEmail (c *gin.Context) {
    token := c.Query("token")
    fmt.Println(token)
    e.db.VerfUserEmail(token)
}
