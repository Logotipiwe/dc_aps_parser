package pkg

import (
	"github.com/gin-gonic/gin"
)

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

func RemoveElement[T comparable](slice []T, element T) []T {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
