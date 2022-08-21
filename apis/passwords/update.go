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

func UpdatePassword(c *gin.Context) {
	var req updatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
		return
	}
	db, err := ginutils.Database(c)
	if err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, err)
		return
	}
	hashedPwd, err := crypto.Encrypt(c, req.Password)
	if err != nil {
		c.Error(err)
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, fmt.Errorf("invalid cipher key"))
		return
	}
	password := &storage.Password{UserEmail: auth.GetCredentials(c).Email, Site: req.Site, SiteUserName: req.SiteUserName}
	if err := db.Model(&password).Updates(storage.Password{HashedPassword: hashedPwd}).Error; err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, fmt.Errorf("error writing user to db: %w", err))
		return
	}
	c.JSON(http.StatusOK, req)
}

