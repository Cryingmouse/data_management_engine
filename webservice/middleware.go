package webservice

import (
	"bytes"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
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
		common.AuditLogger.WithFields(log.Fields{
			"TraceID":     traceID,
			"Method":      c.Request.Method,
			"IP":          c.ClientIP(),
			"URL":         c.Request.URL.Path,
			"RequestBody": common.MaskPassword(string(body)),
		}).Info("Request received")

		common.Logger.WithFields(log.Fields{
			"TraceID":     traceID,
			"Method":      c.Request.Method,
			"IP":          c.ClientIP(),
			"URL":         c.Request.URL.Path,
			"RequestBody": common.MaskPassword(string(body)),
		}).Info("Request received")

		// Process the request
		c.Next()

		// Capture the response body
		responseBody := writer.body.Bytes()

		// Log the response information
		common.AuditLogger.WithFields(log.Fields{
			"TraceID":      traceID,
			"URL":          c.Request.URL.Path,
			"Status":       c.Writer.Status(),
			"ResponseBody": common.MaskPassword(string(responseBody)),
		}).Info("Response")

		common.Logger.WithFields(log.Fields{
			"TraceID":      traceID,
			"URL":          c.Request.URL.Path,
			"Status":       c.Writer.Status(),
			"ResponseBody": common.MaskPassword(string(responseBody)),
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

// I18NTranslatorMiddleware 是i18N验证错误翻译中间件
func I18NTranslatorMiddleware(validate *validator.Validate) gin.HandlerFunc {
	return func(c *gin.Context) {
		en := en_US.New()
		zh := zh_Hans_CN.New()
		UniTranslator = ut.New(en, en, zh)

		USTranslator, _ := UniTranslator.GetTranslator("en_US")
		ZHTranslator, _ := UniTranslator.GetTranslator("zh_Hans_CN")

		en_translations.RegisterDefaultTranslations(validate, USTranslator)
		zh_translations.RegisterDefaultTranslations(validate, ZHTranslator)

		validate.RegisterTranslation("validatePassword", USTranslator, func(ut ut.Translator) error {
			return ut.Add("validatePassword", "Invalid Password.", false) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("validatePassword", fe.Field())

			return t
		})

		validate.RegisterTranslation("validatePassword", ZHTranslator, func(ut ut.Translator) error {
			return ut.Add("validatePassword", "无效的密码.", false) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("validatePassword", fe.Field())

			return t
		})

		c.Next()
	}
}
