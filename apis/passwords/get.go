package passwords

import (
	"fmt"
	"net/http"

	// "strings"

	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/auth"
	"github.com/suyogsoti/password_manager/crypto"

	// "github.com/suyogsoti/password_manager/crypto"
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

func ListPasswords(c *gin.Context) {
	var req listPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
		return
	}
	db, err := ginutils.Database(c)
	if err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, err)
		return
	}
	var passwords []storage.Password
	if err := db.Where(&storage.Password{UserEmail: auth.GetCredentials(c).Email, Site: req.Site}).Find(&passwords).Error; err != nil {
		ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, fmt.Errorf("err getting passwords: %w", err))
		return
	}
	resp := []listPasswordResponse{}
	for _, password := range passwords {
		pswd, err := crypto.Decrypt(c, password.HashedPassword)
		if err != nil {
			ginutils.SetErrorAndAbort(c, http.StatusInternalServerError, fmt.Errorf("error descrypting password for site %q and username %q", password.Site, password.SiteUserName))
			return
		}
		resp = append(resp, listPasswordResponse{password.Site, password.SiteUserName, pswd})
	}
	c.JSON(http.StatusOK, resp)
}
