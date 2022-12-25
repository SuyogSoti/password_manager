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

type updatePasswordRequest struct {
	Site         string `json:"site" binding:"required"`
	SiteUserName string `json:"site_user_name" binding:"required"`
	Password     string `json:"password" binding:"required,min=3"`
}

func UpsertPassword(c *gin.Context) *ginutils.PasswordManagerError {
	var req updatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return ginutils.NewError(http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
	}
	db, err := ginutils.Database(c)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, err)
	}
	hashedPwd, err := crypto.Encrypt(auth.GetCredentials(c).Password, req.Password)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("invalid cipher key"))
	}
	password := &storage.Password{
		UserEmail:      auth.GetCredentials(c).Email,
		Site:           req.Site,
		SiteUserName:   req.SiteUserName,
		HashedPassword: hashedPwd,
	}
	if err := db.Save(&password).Error; err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("error writing user to db: %w", err))
	}
	c.JSON(http.StatusOK, req)
	return nil
}
