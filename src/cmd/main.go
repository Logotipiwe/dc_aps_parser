package main

import (
	"github.com/gin-gonic/gin"
	. "ports-adapters-study/src/internal/adapters"
	. "ports-adapters-study/src/internal/adapters/input"
	. "ports-adapters-study/src/internal/core/application"
	. "ports-adapters-study/src/internal/core/ports"
)

func InitControllers(app App) *gin.Engine {
	router := gin.Default()
	_ = InputPorts{
		ResultPort: NewResultController(router, app.ResultService),
		ParserPort: NewParserController(router, app.ParserService),
	}
	err := router.Run(":81")
	if err != nil {
		panic(err.Error())
	}

	return router
}

func main() {
	println("Started!")
	app := NewApp(CreateAdapters())
	_ = InitControllers(app)
}
