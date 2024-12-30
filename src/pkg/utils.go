package pkg

import "github.com/gin-gonic/gin"

func WithError(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := f(c)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		}
	}
}
