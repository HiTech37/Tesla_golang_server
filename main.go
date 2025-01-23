package main

import (
	"log"
	"tesla_server/controller"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/.well-known", "./static/.well-known")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://test.moovetrax.com", "http://localhost:3000", "https://moovetrax.com"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},                                                 // Allow these HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},                                      // Allow these headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // Allow cookies and HTTP authentication
		MaxAge:           12 * time.Hour, // Cache preflight requests for 12 hours
	}))

	r.GET("/api/tesla_signup", controller.GetTeslaSigninURI)
	r.GET("/api/requestAuth", controller.RequestAuth)
	r.GET("/api/getAllVehichles", controller.GetAllVehicles)
	r.POST("/api/sendCommand", controller.HandleCommand)
	r.POST("/api/connectServer", controller.ConnectDevice)

	// this apis are for testing
	r.GET("api/test", controller.TestFunc)

	// Listen on all network interfaces (allowing access from any IP)
	if err := r.Run("0.0.0.0:8099"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
