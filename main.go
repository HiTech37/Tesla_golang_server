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
		AllowOrigins:     []string{"https://test.moovetrax.com", "http://localhost:3000", "https://moovetrax.com", "https://fleetapi.moovetrax.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	go controller.KafkaConsumer()
	go controller.CronJobs()

	r.GET("/api/tesla_signup", controller.GetTeslaSigninURI)
	r.GET("/api/requestAuth", controller.RequestAuth)
	r.GET("/api/getAllVehichles", controller.GetAllVehicles)
	r.POST("/api/sendCommand", controller.HandleCommand)
	r.POST("/api/connectServer", controller.ConnectDeviceforTest)
	r.POST("/api/getConfigStatus", controller.GetDeviceConfigStatus)
	r.POST("/api/getFleetTelemetryError", controller.GetFleetTelemetryError)
	r.POST("/api/getFleetStatus", controller.GetFleetStatus)
	r.POST("/api/getDeviceLiveData", controller.GetDeviceLiveData)
	r.POST("/api/updateDeviceInfo", controller.UpdateDeviceInfo)

	// this apis are for testing
	r.GET("api/test", controller.TestFunc)

	// Listen on all network interfaces (allowing access from any IP)
	if err := r.Run("0.0.0.0:8099"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
