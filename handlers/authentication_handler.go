package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olefine/quote-id-mock/domain"
)

func generateToken() string {
	width := 128
	buffer := make([]byte, width)
	rand.Read(buffer)

	return base64.RawURLEncoding.EncodeToString(buffer)
}

// HandleAuth handles initial domain.AuthenticationRequest
func HandleAuth(c *gin.Context) {
	var authRequest domain.AuthenticationRequest
	err := c.BindJSON(&authRequest)

	if err != nil {
		log.Fatalln(err.Error())
	}

	db := c.MustGet("tokenRepository").(*domain.TokenRepository)

	if authRequest.IsValid() {
		token := generateToken()
		err = db.Create(token)

		if err != nil {
			errorResp := domain.AuthenticationBadParameters{Error: "Unable to generate token"}
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResp)
		} else {
			resp := domain.NewAuthenticationResponse(token)
			c.JSON(http.StatusCreated, resp)
		}
	} else {
		errorResp := domain.AuthenticationBadParameters{Error: "Not supported params"}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResp)
	}
}
