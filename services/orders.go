package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

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

		// Convert hashed passcode to string
		hashedPasscodeString := string(hashedPasscode)

		// TODO: Create Order and OrderDetail
		order := models.Order{
			ID:         uuid.New().String(),
			Email:      checkoutOrder.Email,
			Address:    checkoutOrder.Address,
			GrandTotal: 0,
			Passcode:   &hashedPasscodeString,
		}
		// Calculate grand total and create detail order
		var detailOrders []models.OrderDetail
		// Loop through products and create detail order
		for _, p := range products {
			detailOrder := models.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQty[p.ID],
				Price:     p.Price,
				Total:     p.Price * int64(orderQty[p.ID]),
			}

			// Calculate grand total
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

		// Don't show passcode in response & only show passcode when order is created
		order.Passcode = &passcode

		// Response with order and detail order
		response := models.OrderWithDetail{
			Order:   order,
			Details: detailOrders,
		}

		// Response with status 201 Created
		ctx.JSON(201, gin.H{
			"message": "Order created successfully",
			"data":    response,
		})
	}
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get id order from url
		orderId := ctx.Param("id")

		// Get the request body
		var confirm models.OrderConfirmation
		if err := ctx.BindJSON(&confirm); err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(400, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Get order from database
		order, err := models.SelectOrderById(db, orderId)
		// Check if error
		if err != nil {
			log.Printf("Error: %v", err)

			// Check if order not found
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(404, gin.H{
					"error": "Order not found",
				})
				return
			}
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
		}

		// Make sure passcode not empty
		if order.Passcode == nil {
			ctx.JSON(400, gin.H{
				"error": "Passcode not found",
			})
			return
		}

		// Compare passcode
		err = bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(confirm.Passcode))
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid passcode",
			})
			return
		}

		// Check if order already paid
		if order.PaidAt != nil {
			ctx.JSON(400, gin.H{
				"error": "Order already paid",
			})
			return
		}

		// Check if grand total and amount not match
		if order.GrandTotal != confirm.Amount {
			ctx.JSON(400, gin.H{
				"error": "Amount not match",
			})
			return
		}

		// TODO: Update order data
		currentTime := time.Now()
		if err := models.UpdateOrderStatus(db, orderId, confirm, currentTime); err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}
		// Formatted time to string
		formattedTime := currentTime.Format("2006-01-02 15:04:05")
		// Don't show passcode in response & only show passcode when order is created
		order.Passcode = nil
		// Update response with paid_at, paid_bank, and paid_account_number
		order.PaidAt = &formattedTime
		order.PaidBank = &confirm.Bank
		order.PaidAccountNumber = &confirm.AccountNumber

		// Response with order data
		ctx.JSON(200, gin.H{
			"message": "Order confirmed successfully",
			"data":    order,
		})

	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get order id from url
		orderId := ctx.Param("id")

		//	Get parameter from query parameter
		//	Example: /orders?id=1
		passcode := ctx.Query("passcode")

		// Get order from database
		order, err := models.SelectOrderById(db, orderId)
		// Check if error
		if err != nil {
			log.Printf("Error: %v", err)

			// Check if order not found
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(404, gin.H{
					"error": "Order not found",
				})
				return
			}
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
		}

		// Make sure passcode not empty
		if order.Passcode == nil {
			ctx.JSON(400, gin.H{
				"error": "Passcode not found",
			})
			return
		}

		// Compare passcode
		err = bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(passcode))
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid passcode",
			})
			return
		}

		// Get order detail from database
		orderDetails, err := models.SelectOrderDetailByOrderId(db, orderId)
		if err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(500, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// Don't show passcode in response & only show passcode when order is created
		order.Passcode = nil

		response := models.OrderWithDetail{
			Order:   order,
			Details: orderDetails,
		}

		// Response with order and detail order
		ctx.JSON(200, gin.H{
			"message": "Order found",
			"data":    response,
		})
	}
}
