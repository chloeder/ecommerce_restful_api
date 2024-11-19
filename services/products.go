package services

import (
	"database/sql"
	"errors"

	"commerce-project/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetProducts(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := models.SelectAllProducts(db)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Failed to get all products",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Get all products",
			"data":    data,
		})
	}
}

func GetProductById(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productId := ctx.Param("id")

		data, err := models.SelectProductById(db, productId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(404, gin.H{
					"message": "Product not found",
				})
			} else {
				ctx.JSON(500, gin.H{
					"message": "Something went wrong",
				})
			}
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Get product by id",
			"data":    data,
		})
	}
}

func CreateProducts(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product models.Product
		err := ctx.BindJSON(&product)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "Invalid request",
			})
			return
		}

		product.ID = uuid.New().String()

		err = models.InsertProduct(db, product)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Failed to create product",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Created product Successfully",
		})
	}
}

func UpdateProducts(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get product id from url
		productId := ctx.Param("id")

		// Data from request body
		var product models.Product
		err := ctx.BindJSON(&product)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "Invalid request",
			})
			return
		}

		// Data from database
		existingProduct, err := models.SelectProductById(db, productId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(404, gin.H{
					"message": "Product not found",
				})
			} else {
				ctx.JSON(500, gin.H{
					"message": "Something went wrong",
				})
			}
			return
		}

		// Check if data from request body is empty
		if product.Name != "" {
			existingProduct.Name = product.Name
		}

		// Check if data from request body is empty
		if product.Price != 0 {
			existingProduct.Price = product.Price
		}

		// Update product data
		err = models.UpdateProduct(db, existingProduct)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Failed to update product",
				"error":   err.Error(),
			})
			return
		}

		// Return updated product data
		ctx.JSON(200, gin.H{
			"message": "Updated product Successfully",
			"data":    product,
		})
	}
}

func SoftDeletedProducts(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get product id from url
		productId := ctx.Param("id")

		// Data from database
		existingProduct, err := models.SelectProductById(db, productId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(404, gin.H{
					"message": "Product not found",
				})
			} else {
				ctx.JSON(500, gin.H{
					"message": "Something went wrong",
				})
			}
			return
		}

		// Soft delete product
		err = models.SoftDeletedProduct(db, existingProduct)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Failed to delete product",
				"error":   err.Error(),
			})
			return
		}

		// Return success message
		ctx.JSON(200, gin.H{
			"message": "Deleted product Successfully",
		})
	}
}

func HardDeletedProducts(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get product id from url
		productId := ctx.Param("id")

		// Data from database
		existingProduct, err := models.SelectProductById(db, productId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(404, gin.H{
					"message": "Product not found",
				})
			} else {
				ctx.JSON(500, gin.H{
					"message": "Something went wrong",
				})
			}
			return
		}

		// Soft delete product
		err = models.HardDeletedProduct(db, existingProduct)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Failed to delete product",
				"error":   err.Error(),
			})
			return
		}

		// Return success message
		ctx.JSON(200, gin.H{
			"message": "Deleted product Successfully",
		})
	}
}
