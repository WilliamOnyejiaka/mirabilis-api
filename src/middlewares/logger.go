package middlewares

import (
	// "log"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"time"
	"os"
)

// func Logger() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		c.Next()
// 		duration := time.Since(start)
// 		log.Printf("%s %s %v", c.Request.Method, c.Request.URL.Path, duration)
// 	}
// }

type ginBodyLogger struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (g *ginBodyLogger) Write(b []byte) (int, error) {
	g.body.Write(b)
	return g.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.InfoLevel)
		// Capture request body
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyString := string(bodyBytes)

		// Create custom response writer
		ginBodyLogger := &ginBodyLogger{ResponseWriter: c.Writer}
		c.Writer = ginBodyLogger

		// Process request
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// Colorize status code
		var statusColor string
		switch {
		case c.Writer.Status() >= 200 && c.Writer.Status() < 300:
			statusColor = "\033[32m" // Green
		case c.Writer.Status() >= 400 && c.Writer.Status() < 500:
			statusColor = "\033[33m" // Yellow
		case c.Writer.Status() >= 500:
			statusColor = "\033[31m" // Red
		default:
			statusColor = "\033[0m" // Reset
		}
		resetColor := "\033[0m"

		// Log with colorized output
		logger.WithFields(logrus.Fields{
			"status":     fmt.Sprintf("%s%d%s", statusColor, c.Writer.Status(), resetColor),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"req_body":   bodyString,
			// "res_body":   ginBodyLogger.body.String(),
			"latency":    latency,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			// "error":      c.Errors.String(),
		}).Info("Request details")
	}
}
