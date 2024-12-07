package songs

import "github.com/gin-gonic/gin"

// Songs HTTP Handlers interface
type Handlers interface {
	Create() gin.HandlerFunc
	Hello() gin.HandlerFunc
}
