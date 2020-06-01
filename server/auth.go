package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

const (
	// The name of the cookie on the client's device.
	CookieName   = "auth_jwt"
	CookiePath   = "/"
	CookieDomain = "127.0.0.1"
	// The name of the service, given in the JWT.
	ServiceName = "lcs-sm"
	// How long, in seconds, a token will stay valid.
	TokenLifetime = 60
	// If a request is made within this many seconds of the token expiring, a new
	// token will automatically be generated and sent.
	TokenRefreshWindow = 20
)

// Reads the secret file for the given encoding method into a byte array that
// can be passed straight into the jwt method for signing tokens. This secret
// can be authomatically generated using a command from the makefile.
func GetSecret(method string) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("jwt_secret.%s", method))
}

// User is used to contain information about the current user.
type User struct {
	Id       int
	Username string
	Role     string
}

// AuthClaims is a model used to encode and parse the claims from the JWT, it
// uses primarily the standard claims with one addition for the Username of the
// current user.
type AuthClaims struct {
	Username string `json:"usr"`
	Role     string `json:"rol"`
	jwt.StandardClaims
}

// GetAuthClaims takes the username of a user and returns an AuthClaims object
// that can be used to generate a JWT.
func GetAuthClaims(user, role string) AuthClaims {
	now := time.Now().Unix()
	return AuthClaims{
		user,
		role,
		jwt.StandardClaims{
			Issuer:    ServiceName,
			IssuedAt:  now,
			ExpiresAt: now + TokenLifetime,
		},
	}
}

// Authorize is middleware used by Gin. Any route that uses the middleware
// will look for a JWT in the request header, parse the claims, and check the
// signature. It will then set flags called "authorization" and "user" based
// on the validity of the given token.
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().Unix()
		cookie, err := c.Cookie(CookieName)

		if err != nil { // No token has been set
			setAnonymous(c)
		} else {
			token, err := jwt.ParseWithClaims(cookie, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
				return GetSecret("hmac")
			})
			if err != nil { // Token is expired or could not be parsed
				setAnonymous(c)
			} else {
				if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid { // Good token
					//log.Println(fmt.Sprintf("Token expires in %d s", claims.ExpiresAt - now))
					if claims.ExpiresAt-now <= TokenRefreshWindow { // Issue new token
						//log.Println("need to make a new token")
						Logout(c)
						Login(c, claims.Username, claims.Role)
					}

					setAuthorized(c, claims.Username, claims.Role)
				} else { // Tokens claims could not be validated
					setAnonymous(c)
				}
			}
		}
	}
}

// CheckAuthenticated ensures that the user has a valid role set by
// Authorize middleware. Otherwise it returns false.
func CheckAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, set := c.Get("role"); !set || role == "" {
			c.Set("auth", false)
		} else {
			c.Set("auth", true)
		}
	}
}

// CheckAdmin ensures that the user has the admin role set by
// Authorize middleware. Otherwise it returns false.
func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, set := c.Get("role"); !set || role != "admin" {
			c.Set("auth", false)
		} else {
			c.Set("auth", true)
		}
	}
}

// HasPermission returns true if the auth flag has been set to
// approved by the authentication middleware. Else it sets an
// unauthorized status on the context and returns false.
func HasPermission(c *gin.Context) bool {
	if auth, exists := c.Get("auth"); !auth.(bool) || !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you are not authorzied to make this request",
		})
		return false
	}
	return true
}

// setAnonymous sets flags to indicate to controller method handlers that
// the user who made the request has not supplied a valid authentication
// token.
func setAnonymous(c *gin.Context) {
	c.Set("role", "")
	c.Set("user", "")
}

// setAuthorized sets flags to indicate to controller method handlers that
// the user who made the request is authorized, and what their username is.
func setAuthorized(c *gin.Context, username, role string) {
	c.Set("role", role)
	c.Set("user", username)
}

// Login takes a username to generate a JWT, and sets that token as a cookie on
// the client. This token will be sent by the client on all future requests and
// will be used to validate the requests by the Authorization middleware. This
// function merely creates a token for a given user. It does not validate that
// the user has supplied a valid password.
func Login(c *gin.Context, user, role string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, GetAuthClaims(user, role))

	secret, err := GetSecret("hmac")
	if err != nil { // Secret file is not found
		panic(err)
	}

	signedTokenString, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}

	c.SetCookie(CookieName, signedTokenString, TokenLifetime, CookiePath, CookieDomain, false, false)
}

// Replaces the JWT on the client's device with an empty cookie.
func Logout(c *gin.Context) {
	c.SetCookie(CookieName, "EXPIRED", TokenLifetime, CookiePath, CookieDomain, false, false)
}

// IsValidAuth takes a database connection, username, and password; and
// determines if the password is valid. Returns true if the user exists and
// the password matches. Returns false if the user could not be found or the
// password does not match.
func IsValidAuth(c *gin.Context, store redis.Conn, user, pass string) bool {
	uid, err := redis.Int(store.Do("GET", fmt.Sprintf("username:%s", user)))
	if err != nil {
		return false
	}
	db_pass, err := redis.String(store.Do("HGET", fmt.Sprintf("user:%d", uid), "password"))
	if err != nil {
		return false
	} else {
		return db_pass == pass
	}
}

// GetUser takes a redis store connection and a username and returns
// a User object corresponding to the username. Returns and error if
// there is an issue connecting to the store or any of the fields
// are not found.
func GetUser(store redis.Conn, username string) (User, error) {
	uid, err := redis.Int(store.Do("GET", fmt.Sprintf("username:%s", username)))
	if err != nil {
		return User{}, err
	}
	userkey := fmt.Sprintf("user:%d", uid)
	if err != nil {
		return User{}, err
	}
	role, err := redis.String(store.Do("HGET", userkey, "role"))
	if err != nil {
		return User{}, err
	}
	return User{uid, username, role}, nil
}
