package main

import (
	"database/sql"
	. "dc-aps-parser/src/internal/adapters/input"
	"dc-aps-parser/src/internal/adapters/output"
	. "dc-aps-parser/src/internal/core/application"
	. "dc-aps-parser/src/internal/core/ports"
	"dc-aps-parser/src/internal/infrastructure"
	"dc-aps-parser/src/internal/infrastructure/tg"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	. "github.com/logotipiwe/dc_go_config_lib"

	"log"
	"os"
)

func main() {
	LoadDcConfig()

	config := infrastructure.NewConfig()
	botAPI := tg.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))

	//db, err := sql.Open("pos", "store.db")
	connectionStr := fmt.Sprintf("postgres://%v:%v@%v:5432/%v?sslmode=disable",
		GetConfig("DB_LOGIN"), GetConfig("DB_PASS"),
		GetConfig("DB_HOST"), GetConfig("DB_NAME"))
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	adapters := OutputPorts{
		TargetClientPort:   output.NewTargetClientWebAdapter(),
		NotificationPort:   output.NewNotificationAdapterTg(botAPI),
		ParsersStoragePort: output.NewParserStorageAdapterPg(db),
	}

	resultService := NewResultService(adapters.TargetClientPort)
	adminService := NewAdminService(config)
	parserNotificationService := NewParserNotificationService(config, adapters.NotificationPort)
	app := App{
		ResultService: resultService,
		ParserService: NewParserService(config, resultService, parserNotificationService, adapters.ParsersStoragePort),
		AdminService:  adminService,
	}

	router := gin.Default()

	inputAdapters := InputAdapters{
		NewParserAdapterHttp(router, app.ParserService),
		NewParserAdapterTg(botAPI, app.ParserService, adminService, parserNotificationService),
	}

	inputAdapters.ParserAdapterTg.InitListening()

	err = router.Run(":81")
	if err != nil {
		log.Fatal(err)
	}
}

type InputAdapters struct {
	*ParserAdapterHttp
	*ParserAdapterTg
}
