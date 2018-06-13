package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sirupsen/logrus"
	"time"
	"io/ioutil"
	"bytes"
	"fmt"
	"github.com/mds1455975151/cmdb/errors"
)

func NotImplemented(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{
		"status":  http.StatusNotImplemented,
		"message": "Not Implemented!",
	})
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type loggerEntryWithFields interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

// Ginrus returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Requests with error are logged using logrus.Error().
// Requests without error are logged using logrus.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func Ginrus(logger loggerEntryWithFields, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method != "POST" && c.Request.Method != "GET" {
			c.Next()
			return
		}

		// some evil middleware modify this values
		path := c.Request.URL.Path

		// Read the Body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)

			// Restore the io.ReadCloser to its original state
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		start := time.Now()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		fields := logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format(timeFormat),
			"response":   blw.body.String(),
		}

		if bodyBytes != nil {
			fields["request"] = string(bodyBytes)
		}

		entry := logger.WithFields(fields)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			entry.Info("Gin Processed")
		}
	}
}

type CommonResponse struct {
	// in: body
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
}

func QuickReplyError(c *gin.Context, err error) {

	if errorCode, ok := err.(*errors.Type); ok {

		QuickReply(c, errorCode.Code, errorCode.Error())
	}
}

func QuickReply(c *gin.Context, code errors.Code, message string, args ... interface{}) {

	resp := CommonResponse{}
	resp.Body.Code = int64(code)

	resp.Body.Message = fmt.Sprintf(message, args...)

	c.JSON(http.StatusOK, resp.Body)
}
