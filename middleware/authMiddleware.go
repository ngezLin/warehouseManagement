package middleware

import (
	"net/http"
	"strings"
	"warehousemanagement/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memeriksa JWT dan menaruh info user di context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow OPTIONS requests to pass through tanpa auth (CORS preflight)
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Ambil token dari "Bearer <token>"
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.VerifyToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Simpan info user di context Gin
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware untuk membatasi akses berdasarkan role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow OPTIONS requests untuk CORS
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		userRole := roleVal.(string)
		for _, r := range allowedRoles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
	}
}
