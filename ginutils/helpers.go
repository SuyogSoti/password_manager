package ginutils

import (
	"fmt"

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
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func SetErrorAndAbort(c *gin.Context, code int, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(code, PasswordManagerError{code, err.Error()})
}

func SetDatabaseInContext(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set(databaseKey, db)
	}
}
