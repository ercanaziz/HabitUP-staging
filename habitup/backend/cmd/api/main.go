package main

import (
	"log"
	"os"

	"habitup/internal/auth"
	"habitup/internal/cache"
	"habitup/internal/db"
	"habitup/internal/habit"
	"habitup/internal/middleware"
	"habitup/internal/queue"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	cache.Connect()
	queue.Connect()
	queue.StartConsumer()

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Auth rotaları (token gerektirmez)
	api := r.Group("/api")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", auth.Register) // Gereksinim 1
			authRoutes.POST("/login", auth.Login)       // Gereksinim 2
		}

		// Korumalı rotalar
		protected := api.Group("/", middleware.AuthRequired())
		{
			protected.POST("/auth/logout", auth.Logout) // Gereksinim 10

			protected.GET("/habits", habit.List)          // Gereksinim 4
			protected.POST("/habits", habit.Create)       // Gereksinim 3
			protected.PUT("/habits/:id", habit.Update)    // Gereksinim 6
			protected.DELETE("/habits/:id", habit.Delete) // Gereksinim 7

			protected.POST("/habits/:id/check", habit.Check)    // Gereksinim 5
			protected.DELETE("/habits/:id/check", habit.Uncheck) // Gereksinim 8
			protected.GET("/habits/:id/stats", habit.Stats)     // Gereksinim 9
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("HabitUp API başlatılıyor — port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
