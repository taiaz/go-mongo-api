package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-mongo-api/config"
	"go-mongo-api/handlers"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests received",
		},
		[]string{"path", "method", "status"},
	)
)

func init() {
	// Register the custom metric with Prometheus's default registry.
	prometheus.MustRegister(requestCount)
}

func main() {
	envMode := os.Getenv("ENV_MODE")

	if envMode == "local" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		log.Println("Loaded environment variables from .env file")
	} else {
		log.Println("Running in production mode, using environment variables from Kubernetes")
	}

	config.ConnectDB()

	router := gin.Default()

	// Middleware để ghi lại số lượng request
	router.Use(func(c *gin.Context) {
		// Tiến hành request
		c.Next()

		// Cập nhật metric sau khi request hoàn thành
		status := c.Writer.Status()
		requestCount.WithLabelValues(c.FullPath(), c.Request.Method, http.StatusText(status)).Inc()
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API endpoints
	router.GET("/test-connection", handlers.TestConnection)
	router.GET("/get-data", handlers.GetData)

	// Run the application
	router.Run(":8080")
}

//export ENV_MODE=local
//go run main.go
