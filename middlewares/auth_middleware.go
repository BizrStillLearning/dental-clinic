package middlewares

import (
	"e-clinic/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.SetCookie("token", "", -1, "/", "localhost", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func RoleBlockMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.HTML(http.StatusForbidden, "error.html", gin.H{"message": "Akses Ditolak"})
			c.Abort()
			return
		}

		userRole := role.(string)
		isAllowed := false
		for _, r := range allowedRoles {
			if userRole == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"message": "Anda tidak memiliki hak akses untuk membuka halaman ini!",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
