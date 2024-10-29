package services

import "github.com/gin-gonic/gin"

func GetProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query := "SELECT id, name, price FROM products"
		ctx.JSON(200, gin.H{
			"message": "Get all products",
			"data":    query,
		})
	}
}

func GetProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Get product by id",
		})
	}
}

func CreateProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Create product",
		})
	}
}

func UpdateProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Update product",
		})
	}
}

func DeleteProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Delete product",
		})
	}
}
