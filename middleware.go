package main

import (
	"encoding/base64"
	"encoding/hex"
	"example.com/covid-go/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/blake2b"
	"log"
	"net/http"
	"strings"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// before request
		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			respondWithError(http.StatusUnauthorized, "Unauthorized", c)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !authenticateUser(pair[0], pair[1]) {
			respondWithError(http.StatusUnauthorized, "Unauthorized", c)
			return
		}
		log.Println("Logged in by", pair[0])
		c.Set(gin.AuthUserKey, pair[0])
		c.Next()
	}
}

func authenticateUser(username, password string) bool {
	user, err := repository.DB.GetUserByUsername(username)
	if err != nil {
		return false
	}
	hash := blake2b.Sum512([]byte(password))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr == user.PasswordHash
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}
