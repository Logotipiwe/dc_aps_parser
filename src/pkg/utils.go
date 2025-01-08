package pkg

import (
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
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

func OsGetNonEmpty(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Value is empty in environment variable " + key)
	}
	return value
}

func OsGetInt64NonEmpty(key string) int64 {
	valueStr := os.Getenv(key)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		panic("Value in environment variable is not int64: " + key)
	}
	return value
}

func Some[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

func OsGetIntNonEmpty(key string) int {
	valueStr := os.Getenv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic("Value in environment variable is not int: " + key)
	}
	return value
}
