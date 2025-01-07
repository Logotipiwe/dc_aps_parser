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
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if v != element {
			result = append(result, v)
		}
	}
	return result
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
