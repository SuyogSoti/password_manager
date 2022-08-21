package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/apis/passwords"
	"github.com/suyogsoti/password_manager/apis/users"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/ginutils"
	"github.com/suyogsoti/password_manager/storage"
)

func indexPage(c *gin.Context) {
	c.JSON(http.StatusOK, struct{ key string }{"Hello World!"})
}
func authenticatedIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You are authenticated " + auth.GetCredentials(c).Email})
}

func main() {
	dsn := "host=localhost user=suyogsoti dbname=password_manager port=5432 sslmode=disable"
	db, err := storage.SetupDB(dsn)
	if err != nil {
		panic("failed to connect database")
	}

	// Set the router as the default one provided by Gin
	router := gin.Default()
	// TODO(suyogsoti): trusting all proxies can be dangerous?
	router.SetTrustedProxies(nil)
	router.Use(ginutils.SetDatabaseInContext(db))

	router.GET("/", indexPage)
	router.POST("/createUser", users.CreateUser)
	router.POST("/authenticate", auth.Authenticate)

	secure := router.Group("/secure")
	{
		secure.Use(auth.CheckAuthenticated())
		secure.GET("/", authenticatedIndex)
		secure.POST("/createPassword", passwords.CreatePassword)
		secure.POST("/listPasswords", passwords.ListPasswords)
	}

	// Start serving the application
	router.Run("localhost:8080")
}
