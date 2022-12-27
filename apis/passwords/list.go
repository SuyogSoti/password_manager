package passwords

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/crypto"
	"github.com/suyogsoti/password_manager/ginutils"
	"github.com/suyogsoti/password_manager/storage"
)

type listPasswordRequest struct {
	Site string `json:"site"`
}

type listPasswordResponse struct {
	Site         string `json:"site" binding:"required"`
	SiteUserName string `json:"site_user_name" binding:"required"`
	Password     string `json:"password" binding:"required,min=3"`
}

func ListPasswords(c *gin.Context) *ginutils.PasswordManagerError {
	var req listPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return ginutils.NewError(http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
	}
	db, err := ginutils.Database(c)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, err)
	}
	var user storage.User
	if err := db.First(&user, auth.GetCredentials(c).Email).Error; err != nil {
		return ginutils.NewError(http.StatusUnauthorized, fmt.Errorf("error getting user: %w", err))
	}
	encryptionKey, err := crypto.Decrypt(auth.GetCredentials(c).Password, user.EncryptedKey)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("error fetching encryption key: %w", err))
	}
	var passwords []storage.Password
	if err := db.Where(&storage.Password{UserEmail: auth.GetCredentials(c).Email, Site: req.Site}).Find(&passwords).Error; err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("err getting passwords: %w", err))
	}
	resp := []listPasswordResponse{}
	for _, password := range passwords {
		pswd, err := crypto.Decrypt(encryptionKey, password.HashedPassword)
		if err != nil {
			return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("error descrypting password for site %q and username %q", password.Site, password.SiteUserName))
		}
		resp = append(resp, listPasswordResponse{password.Site, password.SiteUserName, pswd})
	}
	c.JSON(http.StatusOK, resp)
	return nil
}
