package api

import "github.com/gin-gonic/gin"

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func Error(err error) gin.H {
	return gin.H{"error": err.Error()}
}
