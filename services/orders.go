package services

import (
	"database/sql"
	"log"

	"commerce-project/models"
	"commerce-project/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the request body
		var checkoutOrder models.Checkout
		err := ctx.BindJSON(&checkoutOrder)
		if err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(400, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// List of product IDs and Quantity
		var Ids []string
		orderQty := map[string]int32{}
		for _, o := range checkoutOrder.Products {
			Ids = append(Ids, o.ID)
			orderQty[o.ID] = o.Quantity
		}

		// Get the products from the database
		products, err := models.SelectProductsIn(db, Ids)
		if err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// Check if the product is not found
		if len(products) != len(checkoutOrder.Products) {
			ctx.JSON(400, gin.H{"error": "Product not found"})
			return
		}

		// TODO: Make Passcode and Hashing it
		passcode := utils.PasscodeGenerator(5)
		hashedPasscode, err := bcrypt.GenerateFromPassword([]byte(passcode), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		hashedPasscodeString := string(hashedPasscode)

		// TODO: Create Order and OrderDetail
		order := models.Order{
			ID:         uuid.New().String(),
			Email:      checkoutOrder.Email,
			Address:    checkoutOrder.Address,
			GrandTotal: 0,
			Passcode:   &hashedPasscodeString,
		}

		var detailOrders []models.OrderDetail
		for _, p := range products {
			detailOrder := models.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQty[p.ID],
				Price:     p.Price,
				Total:     p.Price * int64(orderQty[p.ID]),
			}

			order.GrandTotal += detailOrder.Total
			detailOrders = append(detailOrders, detailOrder)
		}

		// Save data order and detail order to database
		err = models.CreateOrder(db, order, detailOrders)
		if err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		order.Passcode = &passcode

		response := models.OrderWithDetail{
			Order:   order,
			Details: detailOrders,
		}

		ctx.JSON(201, gin.H{
			"message": "Order created successfully",
			"data":    response,
		})
	}
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get id order from url

		// Get passcode from request body
	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// code here
	}
}
