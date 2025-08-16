package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	if userType != role {
		return fmt.Errorf("user type does not match expected role")
	}
	return nil
}

func MatchUserTypeToUid(c *gin.Context, userID string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")

	if (userType == "USER" || userType == "user") && uid != userID {
		return fmt.Errorf("user type does not match user ID")
	}

	return CheckUserType(c, userType)
}
