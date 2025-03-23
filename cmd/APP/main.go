package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	initDB()
	router := gin.Default()
	router.POST("/login", login)

	admin := router.Group("/admin")
	admin.Use(authMiddleware("admin"))
	{
		admin.GET("/orders", getOrders)
		admin.POST("/orders", createOrder)
		admin.GET("/orders/:id", getOrderByID)
		admin.PUT("/orders/:id", updateOrder)
		admin.DELETE("/orders/:id", deleteOrder)
	}

	user := router.Group("/user")
	user.Use(authMiddleware("user"))
	{
		user.GET("/orders/:id", getOrderByID)
	}

	router.Run(":8080")
}
