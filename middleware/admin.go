package middleware

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Admin key
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
		key := os.Getenv("ADMINISTRATOR_KEY")

		// Get authorization header
		authorization := ctx.Request.Header.Get("Authorization")

		// Check if authorization header is empty
		if authorization == "" {
			ctx.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		// Check if authorization header is valid
		if authorization != key {
			ctx.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		//	Continue to the next middleware
		ctx.Next()
	}
}
