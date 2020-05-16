package main

import(
  "crypto/sha256"
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

// AuthController holds a redis connection instance and methods
// used to handle user requests related to authentication.
type AuthController struct {
  store redis.Conn
}

// InitAuthController takes a redis connection opject, instantiates
// a new AuthController, and returns a pointer to that controller.
func InitAuthController(store redis.Conn) *AuthController {
  return &AuthController{store}
}

// LoginReqeust is used to hold request header data for login.
type LoginRequest struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

// hash takes a string and returns the sha256 hash of the string.
func hash(value string) string {
  hash := sha256.Sum256([]byte(value))
  return fmt.Sprintf("%x", hash)
}

// Signin takes a request with a json header containing a username
// and password to authenticate a user. It uses the IsValidAuth 
// function to determine if the username and password provided are
// valid, and if they are it sets a cookie on the users device to 
// authenticate next requests and returns an OK status message. 
//
// If the username is not found, or the password is incorrect, a
// BadRequest status is returned.
func (a *AuthController) Signin(c *gin.Context) {
  var request LoginRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  request.Password = hash(request.Password)

  if IsValidAuth(c, a.store, request.Username, request.Password) {
    user, err := GetUser(a.store, request.Username)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": err.Error(),
      })
    }
    Login(c, user.Username, user.Role)
    c.JSON(http.StatusOK, gin.H{
      "success": fmt.Sprintf("authenticated user %v", request.Username),
    })
    return
  } else {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": "incorrect username/password",
    })
    return
  }
}

// Signup takes a request with a json header containing a username 
// and password to create a new user in the data store. It then
// sets a cookie on the users device to authenticate their requests
// and identify them and returns an OK response.
//
// Returns an InternalServerError response if an issue occurs with
// the redis store. 
//
// Returns a bad request if username or password is missing from the
// JSON header object, or if the requested username is already in
// use.
func (a *AuthController) Signup(c *gin.Context) {
  // unmarshall request into LoginRequest
  var request LoginRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  request.Password = hash(request.Password)

  // set they key that will be used to look up a user by username
  usernamekey := fmt.Sprintf("username:%v", request.Username)

  // check to see if username has been used
  exists, err := redis.Bool(a.store.Do("EXISTS", usernamekey))
  if exists {
    c.JSON(http.StatusBadRequest, gin.H{"error": "username is already in use"})
    return
  }
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  // increment the number of users and grab the enxt id
  uid, err := redis.Int(a.store.Do("INCR", "next-user-id"))
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  // set user information in the store
  userkey := fmt.Sprintf("user:%d", uid)
  a.store.Do("HMSET", userkey, "username", request.Username, "password", request.Password, "role", "user")
  a.store.Do("SET", usernamekey, uid)

  // set cookie that authenticates user requests
  Login(c, request.Username, "user")

  // send a success response
  c.JSON(http.StatusOK, gin.H{
    "success": fmt.Sprintf("created user %v", request.Username),
  })
}

// Signout empties the user's authentication token and returns an 
// OK response.
func (a *AuthController) Signout(c *gin.Context) {
  if(HasPermission(c)) {
    Logout(c)

    c.JSON(http.StatusOK, gin.H{"success": "logged out"})
  }
}
