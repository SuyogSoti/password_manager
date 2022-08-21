package users

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/ginutils"
	"github.com/suyogsoti/password_manager/storage"
)

type createUserRequest struct {
	sanitizedUser
	Password string `json:"password" binding:"required,min=3"`
}
type sanitizedUser struct {
	Email string `json:"email" binding:"required,email"`
}

func CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
		return
	}
	db, err := ginutils.Database(c)
	if err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, err)
		return
	}
	hashedPswd, err := auth.HashPassword(req.Password)
	if err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, fmt.Errorf("password can not be hashed: %w", err))
		return
	}
	user := &storage.User{Email: req.Email, HashedPassword: hashedPswd}
	if err := db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			ginutils.SetErrorAndAbort(c, http.StatusBadRequest, fmt.Errorf("user %q already exists", req.Email))
			return
		}
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, fmt.Errorf("error writing user to db: %w", err))
		return
	}
	c.JSON(http.StatusOK, req.sanitizedUser)
}
