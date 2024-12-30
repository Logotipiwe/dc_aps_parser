package driverport

import "github.com/gin-gonic/gin"

type ResultController interface {
	GetResult(ctx *gin.Context) error
	GetResultsHistory(ctx *gin.Context) error
}

type ParserController interface {
	NewParser(ctx *gin.Context) error
	StopParser(ctx *gin.Context) error
}
