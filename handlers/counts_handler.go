package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olefine/quote-id-mock/domain"
)

// HandleCounts checks if token is exists and still valid:
// if token found, valid and body has content: returns payload read from request body
// if token not found or already expiried: returns not_found
// if request body empty: returns bad_request
func HandleCounts(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")

	var reqBody []byte
	var err error
	if reqBody, err = ioutil.ReadAll(c.Request.Body); err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}

	var response *domain.CountsResponse

	// TODO validate if reqBody has content

	rep := c.MustGet("tokenRepository").(*domain.TokenRepository)

	if err = rep.IsTokenValid(authToken); err != nil {
		response = domain.NewCountsResponseError(err.Error())
	} else {
		response = domain.NewCountsResponseSuccess(string(reqBody))
	}

	c.JSON(200, response)
}
