package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/suyogsoti/password_manager/ginutils"
	"github.com/suyogsoti/password_manager/storage"
)

const tokenExpireDuration = time.Hour * 2
const credentialsKey = "credentialsKey"

// Create a struct to read the username and password from the request body
type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type JwtClaims struct {
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=3"`
	jwt.StandardClaims
}

func Authenticate(c *gin.Context) *ginutils.PasswordManagerError {
	var req Credentials
	if err := c.ShouldBindJSON(&req); err != nil {
		return ginutils.NewError(http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
	}
	db, err := ginutils.Database(c)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, err)
	}
	var user storage.User
	if err := db.First(&user, req.Email).Error; err != nil {
		return ginutils.NewError(http.StatusUnauthorized, fmt.Errorf("the email or the user name is incorrect"))
	}
	if !CheckPasswordHash(req.Password, user.HashedPassword) {
		return ginutils.NewError(http.StatusUnauthorized, fmt.Errorf("the email or the user name is incorrect"))
	}
	token, err := GenToken(user.Email, req.Password)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token"))
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
	return nil
}

func CheckAuthenticated(c *gin.Context) *ginutils.PasswordManagerError {
	// There are three ways for the client to carry a Token. 1 Put in request header 2 Put in the request body 3 Put in URI
	// Here, it is assumed that the Token is placed in the Authorization of the Header and starts with Bearer
	// The specific implementation method here should be determined according to your actual business situation
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return ginutils.NewError(http.StatusUnauthorized, fmt.Errorf("request header auth empty"))
	}
	// Split by space
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return ginutils.NewError(http.StatusUnauthorized, fmt.Errorf("request header auth invalid format"))
	}
	// parts[1] is the obtained tokenString. We use the previously defined function to parse JWT to parse it
	mc, err := ParseToken(parts[1])
	if err != nil {
		return ginutils.NewError(http.StatusUnauthorized, fmt.Errorf("invalid token"))
	}
	// Save the currently requested username information to the requested context c
	c.Set(credentialsKey, Credentials{mc.Email, mc.Password})
	c.Next() // Subsequent processing functions can use c.Get("username") to obtain the currently requested user information
	return nil
}

// ParseToken parsing JWT
func ParseToken(tokenString string) (*JwtClaims, error) {
	var claims JwtClaims
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return secrect(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing jwt: %w", err)
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid { // Verification token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// GenToken generates JWT
func GenToken(email, password string) (string, error) {
	// Create our own statement
	c := JwtClaims{
		email,
		password,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),   // Expiration time
			Issuer:    "github.com/suyogsoti/password_manager/auth", // Issuer
		},
	}
	// Creates a signed object using the specified signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// Use the specified secret signature and obtain the complete encoded string token
	return token.SignedString(secrect())
}

func secrect() []byte {
	if secrect := os.Getenv("password_manager_jwt_secrect"); secrect != "" {
		return []byte(secrect)
	}
	return []byte("my secrect key")
}

func GetCredentials(c *gin.Context) Credentials {
	cred, _ := c.Get(credentialsKey)
	return cred.(Credentials)
}
