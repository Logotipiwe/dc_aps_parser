package inputport

import "github.com/gin-gonic/gin"

type ResultPort interface {
	GetResult(ctx *gin.Context) error
	GetResultsHistory(ctx *gin.Context) error
}

type ParserPort interface {
	NewParser(ctx *gin.Context) error
	StopParser(ctx *gin.Context) error
}
