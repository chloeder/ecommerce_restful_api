package routes

import (
	"commerce-project/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Routes() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/products", services.GetProducts())
		v1.GET("/products/{id}", services.GetProductById())
		//v1.POST("/checkout", checkoutProduct)
	}

	admin := router.Group("/api/admin")
	{
		admin.POST("/products", services.CreateProducts())
		admin.PUT("/products/{id}", services.UpdateProducts())
		admin.DELETE("/products/{id}", services.DeleteProducts())
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

}
