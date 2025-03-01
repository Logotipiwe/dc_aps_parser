package main

import (
	"database/sql"
	. "dc-aps-parser/src/internal/adapters"
	. "dc-aps-parser/src/internal/core/application"
	. "dc-aps-parser/src/internal/core/ports"
	"dc-aps-parser/src/internal/infrastructure"
	"dc-aps-parser/src/internal/infrastructure/flyway_pg"
	"dc-aps-parser/src/internal/infrastructure/pg"
	"dc-aps-parser/src/internal/infrastructure/tg"
	"dc-aps-parser/src/internal/infrastructure/web"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	. "github.com/logotipiwe/dc_go_config_lib"

	"log"
)

// TODO упросить получение ссылки
// TODO help админу и юзеру срастить в 1 параметр

func main() {
	LoadDcConfig()

	config := infrastructure.NewConfig()

	connectionStr := fmt.Sprintf("postgres://%v:%v@%v:5432/%v?sslmode=disable",
		GetConfig("DB_LOGIN"), GetConfig("DB_PASS"),
		GetConfig("DB_HOST"), GetConfig("DB_NAME"))
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	fw := flyway_pg.NewFlyway(db, "./data/migrations")
	err = fw.Migrate()
	if err != nil {
		log.Println("Error applying migrations")
		log.Fatal(err)
	}

	botAPI := tg.NewBotAPI(config.TgBotToken)

	resultStorageAdapterPg := pg.NewResultStorageAdapterPg(db)
	adapters := OutputPorts{
		TargetClientPort:       web.NewTargetClientWebAdapter(),
		NotificationPort:       tg.NewNotificationAdapterTg(botAPI),
		ParsersStoragePort:     pg.NewParserStorageAdapterPg(db),
		PermissionsStoragePort: pg.NewPermissionStorageAdapterPg(db),
	}

	resultService := NewResultService(adapters.TargetClientPort)
	adminService := NewAdminService(config)
	parserNotificationService := NewParserNotificationService(config, adapters.NotificationPort)
	permissionsService := NewPermissionsService(config, adapters.PermissionsStoragePort)
	resultStorageService := NewResultsStorageService(resultStorageAdapterPg)
	app := App{
		ResultService: resultService,
		ParserService: NewParserService(
			config,
			resultService,
			parserNotificationService,
			resultStorageService,
			adapters.ParsersStoragePort,
			permissionsService,
		),
		AdminService: adminService,
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
