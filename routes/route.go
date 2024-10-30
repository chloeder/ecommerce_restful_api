package routes

import (
	"commerce-project/services"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Routes(db *sql.DB) error {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/products", services.GetProducts(db))
		v1.GET("/products/:id", services.GetProductById(db))
		//v1.POST("/checkout", checkoutProduct)
	}

	admin := router.Group("/api/admin")
	{
		admin.POST("/products", services.CreateProducts(db))
		admin.PUT("/products/:id", services.UpdateProducts(db))
		admin.DELETE("/products/:id", services.SoftDeletedProducts(db))
		admin.DELETE("/hard-delete/products/:id", services.HardDeletedProducts(db))
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
