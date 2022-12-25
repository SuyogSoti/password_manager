package users

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/crypto"
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

func CreateUser(c *gin.Context) *ginutils.PasswordManagerError {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return ginutils.NewError(http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
	}
	db, err := ginutils.Database(c)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, err)
	}
	hashedPswd, err := auth.HashPassword(req.Password)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("password can not be hashed: %w", err))
	}
	encryptedKey, err := crypto.Encrypt(req.Password, uuid.NewString())
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("unable to create encrypted key: %w", err))
	}
	user := &storage.User{Email: req.Email, HashedPassword: hashedPswd, EncryptedKey: encryptedKey}
	if err := db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return ginutils.NewError(http.StatusBadRequest, fmt.Errorf("user %q already exists", req.Email))
		}
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("error writing user to db: %w", err))
	}
	token, err := auth.GenToken(user.Email, req.Password)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("failed to generate jwt token"))
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
	return nil
}
