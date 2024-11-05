package services

import (
	"commerce-project/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: Get the request body
		var checkoutOrder models.Checkout
		err := ctx.BindJSON(&checkoutOrder)
		if err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(400, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// TODO: Get product from database
		Ids := []string{}
		orderQty := map[string]int32{}
		for _, o := range checkoutOrder.Products {
			Ids = append(Ids, o.ID)
			orderQty[o.ID] = o.Quantity
		}

		products, err := models.SelectProductsIn(db, Ids)
		if err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"data":    products,
			"message": "Checkout success",
		})
	}
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// code here
	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// code here
	}
}
