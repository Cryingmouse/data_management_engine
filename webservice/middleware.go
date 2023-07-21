package webservice

import (
	"bytes"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Middleware function to log request and response information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body
		body, _ := io.ReadAll(c.Request.Body)

		// Restore the request body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Create a custom response writer to capture the response body
		writer := &responseWriterWithCapture{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = writer

		// Get the trace ID from the request context
		traceID := c.Request.Header.Get("X-Trace-ID")

		// Log the request information
		log.WithFields(log.Fields{
			"trace_id":     traceID,
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"ip":           c.ClientIP(),
			"request_body": string(body),
		}).Info("Request received")

		// Process the request
		c.Next()

		// Capture the response body
		responseBody := writer.body.Bytes()

		// Log the response information
		log.WithFields(log.Fields{
			"trace_id":      traceID,
			"status":        c.Writer.Status(),
			"response_body": string(responseBody),
		}).Info("Response")
	}
}

// Custom response writer to capture the response body
type responseWriterWithCapture struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriterWithCapture) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Middleware function to generate and attach a trace ID to the request context
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique trace ID
		traceID := generateTraceID()

		if c.Request.Header.Get("X-Trace-ID") == "" {
			c.Request.Header.Set("X-Trace-ID", traceID)
		}

		// Continue processing the request
		c.Next()
	}
}

// Generate a unique trace ID
func generateTraceID() string {
	// Generate a random number as the trace ID
	seed := time.Now().UnixNano()
	// Use the random generator r to generate random numbers
	randomNumber := rand.New(rand.NewSource(seed)).Intn(999999)

	traceID := strconv.Itoa(randomNumber)

	// Get the current timestamp
	timestamp := time.Now().UnixNano()

	// Combine the random number and timestamp to create a unique trace ID
	traceID = traceID + "_" + strconv.FormatInt(timestamp, 10)

	return traceID
}
