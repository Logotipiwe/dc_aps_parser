package adapters

import (
	inputport "dc-aps-parser/src/internal/core/ports/input"
	"dc-aps-parser/src/pkg"
	"github.com/gin-gonic/gin"
)

type ParserAdapterHttp struct {
	inputport.ParserPort
}

func NewParserAdapterHttp(
	router *gin.Engine,
	service inputport.ParserPort,
) *ParserAdapterHttp {
	p := &ParserAdapterHttp{
		service,
	}
	router.GET("/ping", pkg.WithError(p.Ping))

	return p
}

func (c *ParserAdapterHttp) Ping(ctx *gin.Context) error {
	ctx.JSON(200, nil)
	return nil
}
