package logger

import (
	"bytes"
	"io"
	"log/syslog"
	"os"
	"time"

	"github.com/dedihartono801/api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

func CreateLog(r *gin.Engine) {
	// Create a hook that sends logs to Logstash via TCP
	hook, err := logrus_syslog.NewSyslogHook("tcp", config.GetEnv("LOGSTASH_HOST"), syslog.LOG_INFO, "")
	if err != nil {
		logrus.Error(err.Error())
	} else {
		logrus.AddHook(hook)
	}

	// Create a middleware that logs requests
	r.Use(func(c *gin.Context) {
		// Start timing the request
		start := time.Now()

		// Create a copy of the request body
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		requestBytes, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf) // Reset the request body for further processing

		// Process the request
		c.Next()

		// Log the request
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		logrus.WithFields(logrus.Fields{
			"status":       status,
			"latency":      latency,
			"client_ip":    clientIP,
			"method":       method,
			"path":         path,
			"request_body": string(requestBytes),
		}).Info("Request handled")
	})
}

func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set the log level (e.g., logrus.DebugLevel, logrus.InfoLevel, etc.)
	logrus.SetLevel(logrus.InfoLevel)

	// Optionally, set the output to a file in addition to the standard output.
	file, err := os.OpenFile("logs/activity.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}
}
