package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go-mongo-api/config"
	"go-mongo-api/handlers"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type logrusWriter struct{}

func (w logrusWriter) Write(p []byte) (n int, err error) {
	logMessage := string(p)
	logMessage = strings.TrimSuffix(logMessage, "\n")
	lowerCaseMessage := strings.ToLower(logMessage)

	switch {
	case strings.Contains(lowerCaseMessage, "fatal"):
		logrus.Fatal(logMessage)
	case strings.Contains(lowerCaseMessage, "error"):
		logrus.Error(logMessage)
	case strings.Contains(lowerCaseMessage, "warn"):
		logrus.Warn(logMessage)
	default:
		logrus.Info(logMessage)
	}

	return len(p), nil
}

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
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	log.SetOutput(logrusWriter{})
	log.SetFlags(0)

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

	router.Use(func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()
		requestCount.WithLabelValues(c.FullPath(), c.Request.Method, http.StatusText(status)).Inc()
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API endpoints
	router.GET("/test-connection", handlers.TestConnection)
	router.GET("/get-data", handlers.GetData)

	go func() {
		for {
			log.Println("Logging message every 3 minutes...")
			time.Sleep(3 * time.Minute)
		}
	}()

	// Run the application
	router.Run(":8080")
}

// export ENV_MODE=local
// go run main.go
// log.SetOutput(os.Stdout)
