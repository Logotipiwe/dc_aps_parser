package tests

import (
	"dc-aps-parser/src/internal/core/application"
	. "dc-aps-parser/src/internal/core/domain"
	"dc-aps-parser/src/internal/infrastructure"
	"dc-aps-parser/src/internal/infrastructure/mock"
	"dc-aps-parser/src/pkg"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
	"time"
)

type AdapterMocks struct {
	notification       *mock.NotificationAdapterMock
	targetClient       *mock.TargetClientAdapterMock
	permissionsStorage *mock.PermissionStorageAdapterMock
	parserStorage      *mock.ParserStorageAdapterMock
}

type AppBuilder struct {
	*infrastructure.Config
	AdapterMocks
	*application.App
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func (a *AppBuilder) WithAdapterMocks() *AppBuilder {
	a.AdapterMocks = AdapterMocks{
		mock.NewNotificationAdapterMock(),
		mock.NewTargetClientAdapterMock(),
		mock.NewPermissionStorageAdapterMock(),
		mock.NewParserStorageAdapterMock(),
	}
	return a
}

func (a *AppBuilder) WithConfigMock() *AppBuilder {
	a.Config = &infrastructure.Config{
		ParseInterval:                 time.Millisecond * 100,
		TgAdminChatId:                 int64(10),
		DefaultAllowedApsNum:          200,
		TgParserLaunchMessage:         "Parser launched",
		TgUserStartMessage:            "Hello!",
		TgParserAlreadyStoppedMessage: "Parser already stopped",
		TgErrorStoppingParserMessage:  "Error stopping parser",
		TgParserStoppedMessage:        "Parser stopped",
		TgAdminHelpMessage:            "help admin",
		TgUserHelpMessage:             "help user",
		TgStoppedParserStatusMessage:  "parser is stopped",
		TgUnknownCommandMessage:       "Unknown command",
		TgActiveParserStatus:          "Parser is active",
		TgInitialApsCountFormat:       "init %d aps",
		TgApsNumNotAllowedFormat:      "Yours %d, allowed %d",
		TgErrorMessage:                "Yps, err",
	}
	return a
}

func (a *AppBuilder) WithDefaultMocks() *AppBuilder {
	return a.WithConfigMock().WithAdapterMocks()
}

func (a *AppBuilder) WithParserStorage(parsers []ParserData) *AppBuilder {
	a.parserStorage.SetParsers(parsers)
	return a
}

func (a *AppBuilder) WithClientResults(results []int) *AppBuilder {
	a.targetClient.SetResults(results)
	return a
}

func (a *AppBuilder) Build() *AppBuilder {
	resultService := application.NewResultService(a.AdapterMocks.targetClient)
	parserNotificationService := application.NewParserNotificationService(a.Config, a.AdapterMocks.notification)
	permissionsService := application.NewPermissionsService(a.Config, a.AdapterMocks.permissionsStorage)

	a.App = &application.App{
		ResultService: resultService,
		ParserService: application.NewParserService(
			a.Config,
			resultService,
			parserNotificationService,
			a.AdapterMocks.parserStorage,
			permissionsService,
		),
		AdminService: application.NewAdminService(a.Config),
	}
	return a
}

func initAppWithMocks() (application.App, AdapterMocks, *infrastructure.Config) {
	app := NewAppBuilder().WithDefaultMocks().Build()
	return *app.App, app.AdapterMocks, app.Config
}

func Test_ParserLaunch(t *testing.T) {
	t.Run("Parser creates", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		app, adapterMocks, config := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, 1, len(sentMessages))
		assert.Equal(t, config.TgParserLaunchMessage, sentMessages[0].Text)
		assert.Equal(t, int64(1), sentMessages[0].ChatID)
	})

	t.Run("Parser doesn't stop if not exist", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, _, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()
		assert.Error(t, app.ParserService.StopParser(1))
	})

	t.Run("Many parsers create", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, config := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		adapterMocks.notification.SetCalls(6)

		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.NoError(t, err)
		_, err = app.ParserService.LaunchParser(ParserParams{ChatID: 2, UserName: "username"})
		assert.NoError(t, err)
		_, err = app.ParserService.LaunchParser(ParserParams{ChatID: 3, UserName: "username"})
		assert.NoError(t, err)

		adapterMocks.notification.WaitForCalls()

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, 6, len(sentMessages))

		check := func(chatId int64, messages []mock.SentMessageMock) {
			messagesToChat := pkg.Filter(messages, func(msg mock.SentMessageMock) bool {
				return msg.ChatID == chatId
			})
			assert.Equal(t, 2, len(messagesToChat))
			assert.Equal(t, config.TgParserLaunchMessage, messagesToChat[0].Text)
			assert.Equal(t, fmt.Sprintf(config.TgInitialApsCountFormat, 0), messagesToChat[1].Text)
		}

		check(1, sentMessages)
		check(2, sentMessages)
		check(3, sentMessages)
	})

	t.Run("Parser gets init aps count", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		app, adapterMocks, config := initAppWithMocks()
		adapterMocks.targetClient.SetResults([]int{3})
		defer app.ParserService.StopAllParsersSync()

		adapterMocks.notification.SetCalls(2)
		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		adapterMocks.notification.WaitForCalls()

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, config.TgParserLaunchMessage, sentMessages[0].Text)
		assert.Equal(t, fmt.Sprintf(config.TgInitialApsCountFormat, 3), sentMessages[1].Text)
	})

	t.Run("Parser restarts correctly with new links", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		app, adapterMocks, config := initAppWithMocks()
		adapterMocks.targetClient.SetResults([]int{3})
		defer app.ParserService.StopAllParsersSync()

		adapterMocks.notification.SetCalls(8)
		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}
		_, err = app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}
		_, err = app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}
		parser, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		adapterMocks.notification.WaitForCalls()

		sentMessages := adapterMocks.notification.GetSentMessages()
		parserStartedMessages := pkg.Filter(sentMessages, func(msg mock.SentMessageMock) bool {
			return msg.Text == config.TgParserLaunchMessage
		})
		initialApsCountMessages := pkg.Filter(sentMessages, func(msg mock.SentMessageMock) bool {
			return msg.Text == fmt.Sprintf(config.TgInitialApsCountFormat, 3)
		})
		assert.Equal(t, 4, len(parserStartedMessages))
		assert.Equal(t, 4, len(initialApsCountMessages))
		assert.Equal(t, 1, len(app.ParserService.GetActiveParsers()))
		assert.Equal(t, parser.ID, app.ParserService.GetActiveParsers()[0].ID)
	})
}

func TestParserWorks(t *testing.T) {
	t.Run("Parser detects new ap", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		adapterMocks.targetClient.SetResults([]int{10, 10, 10, 11})
		adapterMocks.notification.SetCalls(3)
		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.NoError(t, err)
		adapterMocks.notification.WaitForCalls()

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, "Новое объявление link_11", sentMessages[2].Text)
	})

	t.Run("Parser detects many new aps", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		adapterMocks.targetClient.SetResults([]int{10, 10, 11, 11, 12})
		adapterMocks.notification.SetCalls(4)
		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.NoError(t, err)
		adapterMocks.notification.WaitForCalls()

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, "Новое объявление link_11", sentMessages[2].Text)
		assert.Equal(t, "Новое объявление link_12", sentMessages[3].Text)
	})

	t.Run("Parser detects many new aps at one time", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		adapterMocks.targetClient.SetResults([]int{10, 14})
		adapterMocks.notification.SetCalls(6)
		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.NoError(t, err)
		adapterMocks.notification.WaitForCalls()

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, 6, len(sentMessages))
	})

	t.Run("Parser ignores hiding and showing some ap", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		adapterMocks.targetClient.SetResults([]int{10, 9, 10, 11})
		adapterMocks.notification.SetCalls(3)
		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.NoError(t, err)
		adapterMocks.notification.WaitForCalls()

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, "Новое объявление link_11", sentMessages[2].Text)
	})
}

func TestParserAutoStart(t *testing.T) {
	t.Run("Nothing starts if storage empty", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		initAppWithMocks()
	})

	t.Run("All stored parsers start", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app := NewAppBuilder().WithDefaultMocks().
			WithParserStorage([]ParserData{
				{1, "1", "username"},
				{2, "2", "username"},
				{3, "3", "username"},
			}).
			WithClientResults([]int{10, 11}).
			Build()
		defer app.ParserService.StopAllParsersSync()

		assert.Equal(t, 3, len(app.ParserService.GetActiveParsers()))
	})
}

func TestParserPermissions(t *testing.T) {
	t.Run("Denies if more than config", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, config := initAppWithMocks()

		apsNum := config.DefaultAllowedApsNum + 1
		adapterMocks.targetClient.SetResults([]int{apsNum})
		parser, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.Error(t, err)
		assert.Nil(t, parser)
		var notAllowedErr NotAllowedError
		assert.ErrorAs(t, err, &notAllowedErr)
		assert.Equal(t, config.DefaultAllowedApsNum, notAllowedErr.AllowedNum)
		assert.Equal(t, apsNum, notAllowedErr.RequestedNum)
	})

	t.Run("Denies if more than config and storage", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, config := initAppWithMocks()

		apsNum := config.DefaultAllowedApsNum + 1
		adapterMocks.targetClient.SetResults([]int{apsNum})
		adapterMocks.permissionsStorage.SetPermissions([]mock.PermissionMock{{1, config.DefaultAllowedApsNum}})
		parser, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.Error(t, err)
		assert.Nil(t, parser)
		var notAllowedErr NotAllowedError
		assert.ErrorAs(t, err, &notAllowedErr)
		assert.Equal(t, config.DefaultAllowedApsNum, notAllowedErr.AllowedNum)
		assert.Equal(t, apsNum, notAllowedErr.RequestedNum)
	})

	t.Run("Allows if storage allows but config doesn't", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, config := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		allowedInStorage := config.DefaultAllowedApsNum + 10
		apsNum := config.DefaultAllowedApsNum + 9
		adapterMocks.permissionsStorage.SetPermissions([]mock.PermissionMock{{1, allowedInStorage}})
		adapterMocks.targetClient.SetResults([]int{apsNum})
		parser, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.NoError(t, err)
		assert.NotNil(t, parser)
	})

	t.Run("Denies if storage denies but config allow", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, config := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		allowedInStorage := config.DefaultAllowedApsNum - 10
		apsNum := config.DefaultAllowedApsNum - 5
		adapterMocks.permissionsStorage.SetPermissions([]mock.PermissionMock{{1, allowedInStorage}})
		adapterMocks.targetClient.SetResults([]int{apsNum})
		parser, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, UserName: "username"})
		assert.Error(t, err)
		assert.Nil(t, parser)
		var notAllowedErr NotAllowedError
		assert.ErrorAs(t, err, &notAllowedErr)
		assert.Equal(t, allowedInStorage, notAllowedErr.AllowedNum)
		assert.Equal(t, apsNum, notAllowedErr.RequestedNum)
	})
}
