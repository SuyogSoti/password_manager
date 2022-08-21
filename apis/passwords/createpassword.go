package passwords

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/crypto"
	"github.com/suyogsoti/password_manager/ginutils"
	"github.com/suyogsoti/password_manager/storage"
)

type createPasswordRequest struct {
	Site         string `json:"site" binding:"required"`
	SiteUserName string `json:"site_user_name" binding:"required"`
	Password     string `json:"password" binding:"required,min=3"`
}

func CreatePassword(c *gin.Context) {
	var req createPasswordRequest
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
	password := &storage.Password{UserEmail: auth.GetCredentials(c).Email, Site: req.Site, SiteUserName: req.SiteUserName, HashedPassword: hashedPwd}
	if err := db.Create(password).Error; err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			ginutils.SetErrorAndAbort(c, http.StatusBadRequest, fmt.Errorf("site %q already has a password for user %q", req.Site, req.SiteUserName))
			return
		}
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, fmt.Errorf("error writing user to db: %w", err))
		return
	}
	c.JSON(http.StatusOK, req)
}
