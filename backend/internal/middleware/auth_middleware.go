package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	jwtutil "github.com/smazmi/team-task-manager-assignment/backend/pkg/jwt"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/response"
)

const (
	ContextUserIDKey    = "user_id"
	ContextUserEmailKey = "user_email"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperror.Unauthorized("authorization header is required"))
			c.Abort()
			return
		}

		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || !strings.EqualFold(tokenParts[0], "Bearer") {
			response.Error(c, apperror.Unauthorized("authorization header must use the Bearer scheme"))
			c.Abort()
			return
		}

		claims, err := jwtutil.ParseToken(m.jwtSecret, tokenParts[1])
		if err != nil {
			response.Error(c, apperror.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Next()
	}
}

func GetCurrentUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}
