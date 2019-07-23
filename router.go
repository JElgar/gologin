package main

import (
   models "github.com/jelgar/login/models"
   email "github.com/jelgar/login/email"
   errors "github.com/jelgar/login/errors"
   "fmt"
   "github.com/gin-gonic/gin"
   "time"
   "net/http"
//   "net/url"
)

func SetupRouter(env *Env) *gin.Engine {
    r := gin.Default()
    // guessing this is pretty handy for version control :D
    api := r.Group("api/v1")

    tokenAuth := api.Group("/")
    tokenAuth.Use(AuthRequired())


    api.GET("/ping", ping)
    api.GET("/user", env.getUser)
    api.POST("/createUser", env.createUser)
    api.POST("/login", env.login)
    api.GET("/resetPasswordRequest", env.passResetRequest)
    api.GET("/resetPassword", env.passReset)
    api.POST("/sendMail", env.sendMail)
    api.GET("/confirmEmail", env.confirmEmail)
    tokenAuth.GET("/welcome", welcome)

    return r
}

// TODO do i check that the token hasnt expired
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("Hello")
        cookie, err := c.Cookie("token")
        fmt.Println(cookie)

        //cookie, err := c.Request.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                // If the cookie is not set, return an unauthorized status
                fmt.Println("Cookie not set")
                c.JSON(400, errors.ApiError{err, "No cookie supplied", 400})
	            return
            }
            c.JSON(500, errors.ApiError{err, "Error getting cookie", 400})
	    	return
	    }
        //tokenString, err := url.QueryUnescape(cookie.Value)
        claims := &models.Claims{}

        token, erro := models.ParseWClaims(cookie, claims)

        if !token.Valid {
            fmt.Println("Invalid Token")
            return
        }

        if erro != nil {
            if erro == models.ErrSignatureInvalid{
                fmt.Println("Invalid Signature")
                return
            }
            fmt.Println("Unknown error")
            return
        }
        c.Next()
    }
}

// TODO Hope there is a better way to do this than needin this function havin already authed. Need to look at concurency at some point 
// I could just use an auth function at the start of each secure endpoint but i feel like this is a bit naff
func getClaimsNoErrChecking(c *gin.Context) *models.Claims {
        cookie, _ := c.Cookie("token")
        claims := &models.Claims{}

        models.ParseWClaims(cookie, claims)
        return claims
}

func ping(c *gin.Context){
    c.JSON(200, gin.H{
        "world": "Hello",
    })
}

// This is a function to test the json webtoken authentication
// TODO split this authorization step into an external jwt function
// Or middleware that works too :D 
func welcome (c *gin.Context) {
    claims := getClaimsNoErrChecking(c)
    c.JSON(200, gin.H{
        "Message": "Hello, ",
        "Username": claims.Username,
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
            // TODO Definately dont panic here do c.thingy or whatever
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
        // TODO STOP PANICING
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
        // THis isnt really needed anymore
        //switch err.Code {
        //    // 401 -> Incorrect password
        //    case 401:
        //        fmt.Println("Incorrect Password")
        //        c.JSON(err.Code, err)
        //        return
        //    case 404:
        //        fmt.Println("User does not exist")
        //        c.JSON(err.Code, err)
        //        return 
        //    case 500:
        //        fmt.Println("Uknown error so so sorry")
        //        c.JSON(err.Code, err)
        //        return 
        //    default:
        //        fmt.Println("retunr a 500 -> Unknown error")
        //        c.JSON(err.Code, err)
        //        return 
        //}
        c.JSON(err.Code, err)
        return
    }
    fmt.Println(user)
    // If user exists return a JWT being like yup and err nill
    // Otherwise return no JWT and be like that this was the error -> eg no user

    // TODO Can i do some of this elsewhere or it this alright?

    expirationTime := time.Now().Add(5 * time.Minute)

    claims := &models.Claims {
        Username: user.Username,
        StandardClaims: models.NewStandardClaims(expirationTime),
    }
    token := models.NewJWT(models.DefaultSignMethod(), claims)

    tokenString, erro := token.SignedString(models.GetKey())
    if erro != nil {
        fmt.Println("Error making token into signed string")
    }
   
    // TODO making these both false seemed to fix an issue but i dont want them to both be false im guessing 
    c.SetCookie(
        "token",
        tokenString,
        3600,
        "/",
        "",
        false,
        false)
}

func (e *Env) passResetRequest(c *gin.Context) {
   // Find user in database
    var u models.User
    c.BindJSON(&u)
    //user, err := e.db.GetUser(&models.User{Username: "john", Password:"123"})
    user, err := e.db.GetUser(&u)
    if err != nil {
        c.JSON(err.Code, err)
        return
    }
    fmt.Println(user)
   // Set token in datbases for password rest
    
   // Send email to user's email address with custom url

   // 
}

func (e *Env) passReset(c *gin.Context) {
    //token := c.Query("token")
    //e.db.PasswordReset(token)
}

// ONLY FOR TESTING
func (e *Env) sendMail (c *gin.Context) {
    // This is a test handler to send emails to a user
    err := email.Send("James", "jamezy850@gmail.com", "jameselgar.com", "email/email.html")
    if err != nil {
        panic(err)
    }

}

func (e *Env) confirmEmail (c *gin.Context) {
    token := c.Query("token")
    e.db.VerfUserEmail(token)
}
