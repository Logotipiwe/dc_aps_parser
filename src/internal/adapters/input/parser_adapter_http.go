package input

import (
	inputport "dc-aps-parser/src/internal/core/ports/input"
	"dc-aps-parser/src/pkg"
	"github.com/gin-gonic/gin"
	"strconv"
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
	parserApi := router.Group("/parser")
	parserApi.POST("/new", pkg.WithError(p.NewParser))
	parserApi.POST("/stop", pkg.WithError(p.StopParser))

	return p
}

func (c *ParserAdapterHttp) NewParser(ctx *gin.Context) error {
	parser, err := c.ParserPort.NewParser(0, "")
	if err != nil {
		return err
	}
	ctx.JSON(200, parser)
	return nil
}

func (c *ParserAdapterHttp) StopParser(ctx *gin.Context) error {
	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}
	c.ParserPort.StopParser(id)
	ctx.JSON(200, nil)
	return nil
}
