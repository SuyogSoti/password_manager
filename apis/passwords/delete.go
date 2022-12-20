package passwords

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/ginutils"
	"github.com/suyogsoti/password_manager/storage"
)

type deletePasswordRequest struct {
	Site         string `json:"site" binding:"required"`
	SiteUserName string `json:"site_user_name" binding:"required"`
}

func DeletePassword(c *gin.Context) *ginutils.PasswordManagerError {
	var req deletePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return ginutils.NewError(http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
	}
	db, err := ginutils.Database(c)
	if err != nil {
		return ginutils.NewError(http.StatusInternalServerError, err)
	}
	password := &storage.Password{UserEmail: auth.GetCredentials(c).Email, Site: req.Site, SiteUserName: req.SiteUserName}
	if err := db.Unscoped().Delete(password).Error; err != nil {
		return ginutils.NewError(http.StatusInternalServerError, fmt.Errorf("error deleting password from db: %w", err))
	}
	c.JSON(http.StatusOK, req)
	return nil
}
