package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/apis/passwords"
	"github.com/suyogsoti/password_manager/apis/users"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/ginutils"
	"github.com/suyogsoti/password_manager/storage"
	"gorm.io/driver/postgres"
)

func indexPage(c *gin.Context) {
	c.JSON(http.StatusOK, struct{ key string }{"Hello World!"})
}
func authenticatedIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You are authenticated " + auth.GetCredentials(c).Email})
}

func main() {
	config := postgres.Config{
		DSN: "host=localhost user=suyogsoti dbname=password_manager port=5432 sslmode=disable",
	}
	if dsn := os.Getenv("password_manager_postgres_dsn"); dsn != "" {
		log.Println("using password_manager_postgres_dsn dsn")
		config = postgres.Config{
			DriverName: "cloudsqlpostgres",
			DSN:        dsn,
		}
	}
	db, err := storage.SetupDB(config)
	if err != nil {
		log.Panicf("failed to connect database, %v", err)
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
		secure.POST("/updatePasswords", passwords.UpdatePassword)
	}

	if os.Getenv("password_manager_env") == "prod" {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"https://github.com"},
			AllowMethods:     []string{"POST", "GET"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			// AllowOriginFunc: func(origin string) bool {
			// 	return origin == "https://github.com"
			// },
			MaxAge: 12 * time.Hour,
		}))
		router.Run()
	} else {
		// Start serving the application
		router.Run("localhost:8080")
	}
}
