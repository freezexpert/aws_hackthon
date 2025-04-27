package cors

import (
	"time"

	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Default is a default setting of cors for go-gin
func Default() gin.HandlerFunc {

	corsConfig := ginCors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"POST, GET, OPTIONS, PUT, DELETE, UPDATE"},
		AllowHeaders:     []string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}
	return ginCors.New(corsConfig)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
