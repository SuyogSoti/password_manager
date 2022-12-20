package ginutils

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const databaseKey = "databaseKey"

// Database fetches the current database from context
func Database(c *gin.Context) (*gorm.DB, error) {
	dbAnyType, ok := c.Get(databaseKey)
	if !ok {
		return nil, fmt.Errorf("database not connected")
	}
	return dbAnyType.(*gorm.DB), nil
}

type PasswordManagerError struct {
	Code    int   `json:"code" binding:"required"`
	Message error `json:"message" binding:"required"`
}

func (p PasswordManagerError) Error() string {
	return p.Message.Error()
}

func NewError(code int, message error) *PasswordManagerError {
	return &PasswordManagerError{code, message}
}

func SetErrorAndAbort(c *gin.Context, err *PasswordManagerError) error {
	if err == nil {
		return nil
	}
	c.Error(err)
	c.AbortWithStatusJSON(err.Code, err)
	return err
}

func WrapHandler(handler func(c *gin.Context) *PasswordManagerError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			SetErrorAndAbort(c, err)
		}
	}
}

func SetGet(r *gin.RouterGroup, path string, handler func(c *gin.Context) *PasswordManagerError) {
	r.GET(path, WrapHandler(handler))
}
func SetPost(r *gin.RouterGroup, path string, handler func(c *gin.Context) *PasswordManagerError) {
	r.GET(path, WrapHandler(handler))
}
func SetMiddleWare(r *gin.RouterGroup, handler func(c *gin.Context) *PasswordManagerError) {
	r.Use(WrapHandler(handler))
}

func SetDatabaseInContext(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set(databaseKey, db)
	}
}

func LogError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %v", msg, err)
	}
}
