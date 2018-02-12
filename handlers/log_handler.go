package handlers

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/olefine/quote-id-mock/domain"
)

// HandleLogs show gin logs using HTTP2 streaming
func HandleLogs(c *gin.Context) {
	logger := c.MustGet("loggingWriter").(*domain.HTTP2Writer)

	c.Header("Content-Type", "text/html")
	c.Header("Transfer-Encoding", "chunked")
	c.Stream(func(w io.Writer) bool {
		logEntry, exists := <-logger.Out
		if exists {
			toWrite := append(logEntry.([]byte), []byte("<br>")...)
			w.Write(toWrite)
		}
		return true
	})
}
