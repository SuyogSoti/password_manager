package passwords

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/ginutils"
	"github.com/suyogsoti/password_manager/storage"
)

type deletePasswordRequest struct {
	Site         string `json:"site" binding:"required"`
	SiteUserName string `json:"site_user_name" binding:"required"`
}

func DeletePassword(c *gin.Context) {
	var req deletePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
		return
	}
	db, err := ginutils.Database(c)
	if err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, err)
		return
	}
	password := &storage.Password{UserEmail: auth.GetCredentials(c).Email, Site: req.Site, SiteUserName: req.SiteUserName}
	if err := db.Unscoped().Delete(password).Error; err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, fmt.Errorf("error deleting password from db: %w", err))
		return
	}
	c.JSON(http.StatusOK, req)
}
