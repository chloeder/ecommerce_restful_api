package routes

import (
	"database/sql"
	"net/http"
	"time"

	"commerce-project/middleware"
	"commerce-project/services"
	"github.com/gin-gonic/gin"
)

func Routes(db *sql.DB) error {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/products", services.GetProducts(db))
		v1.GET("/products/:id", services.GetProductById(db))
		v1.POST("/checkout", services.CheckoutOrder(db))
	}

	admin := router.Group("/api/admin")
	{
		admin.POST("/products", middleware.AdminMiddleware(), services.CreateProducts(db))
		admin.PUT("/products/:id", middleware.AdminMiddleware(), services.UpdateProducts(db))
		admin.DELETE("/products/:id", middleware.AdminMiddleware(), services.SoftDeletedProducts(db))
		admin.DELETE("/hard-delete/products/:id", middleware.AdminMiddleware(), services.HardDeletedProducts(db))
	}

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}

	return nil
}
