package main

import (
	"github.com/gin-gonic/gin"
	"ports-adapters-study/src/internal/adapters/input"
	"ports-adapters-study/src/internal/core/application"
	"ports-adapters-study/src/internal/platform/storage"
)

func main() {
	println("Started!")

	handler := input.NewResultHandler(
		*application.NewResultService(
			storage.NewResultRepository(),
		),
	)

	router := gin.Default()
	router.GET("/get", withError(func(c *gin.Context) error {
		err := handler.GetResult(c)
		if err != nil {
			return err
		}
		return nil
	}))
	router.GET("/history", withError(func(c *gin.Context) error {
		err := handler.GetResultsHistory(c)
		if err != nil {
			return err
		}
		return nil
	}))

	err := router.Run(":81")
	if err != nil {
		panic(err.Error())
	}
}

func withError(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := f(c)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		}
	}
}
