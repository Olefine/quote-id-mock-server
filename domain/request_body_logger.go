package domain

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// RequestLogger log request body for post requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		toReturn := ioutil.NopCloser(bytes.NewBuffer(buf))
		toHandleBody := ioutil.NopCloser(bytes.NewBuffer(buf))

		handleLog(toHandleBody)

		c.Request.Body = toReturn
		c.Next()
	}
}

func handleLog(r io.Reader) {
	buf := new(bytes.Buffer)
	read, _ := buf.ReadFrom(r)

	if read > 0 {
		rawLogLine := buf.String()

		formatted := fmt.Sprintf("[APP] Body - %v\n", rawLogLine)
		gin.DefaultWriter.Write([]byte(formatted))
	}
}
