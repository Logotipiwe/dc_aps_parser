package tests

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"ports-adapters-study/src/internal/adapters/output"
	"ports-adapters-study/src/internal/core/application"
	"ports-adapters-study/src/internal/core/domain"
	"ports-adapters-study/src/internal/core/ports"
	"testing"
)

func initAppWithMocks(expectedResults []*domain.ParseResult) (*output.NotificationAdapterMock, application.App) {
	resultStorageMock := output.NewResultStorageMock()
	targetClientAdapterMock := output.NewTargetClientAdapterMock(expectedResults)
	notificationAdapterMock := output.NewNotificationAdapterMock()
	app := application.NewApp(ports.OutputPorts{
		ResultStoragePort: resultStorageMock,
		TargetClientPort:  targetClientAdapterMock,
		NotificationPort:  notificationAdapterMock,
	})
	return notificationAdapterMock, app
}

func Test_CreateParser(t *testing.T) {
	t.Run("Parser creates", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		notificationAdapterMock, app := initAppWithMocks([]*domain.ParseResult{
			{ID: 0, ApsNum: 0},
		})
		defer app.ParserService.StopAllParsersSync()

		parser, err := app.ParserService.NewParser()

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		if parser.ID != 0 {
			t.Errorf("GetResult() parser id = %v, want %v", parser.ID, 0)
		}

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, 1, len(sentMessages))
		assert.Equal(t, "Парсер 0 запущен", sentMessages[0])

	})

	t.Run("Parser works", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		notificationAdapterMock, app := initAppWithMocks([]*domain.ParseResult{
			{ID: 0, ApsNum: 100},
			{ID: 0, ApsNum: 100},
			{ID: 0, ApsNum: 101},
			{ID: 0, ApsNum: 103},
			{ID: 0, ApsNum: 103},
			{ID: 0, ApsNum: 100},
		})
		defer app.ParserService.StopAllParsersSync()

		notificationAdapterMock.SetCalls(4)

		_, err := app.ParserService.NewParser()

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		notificationAdapterMock.WaitForCalls()

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, 4, len(sentMessages))
		assert.Equal(t, "Квартир стало больше на 1", sentMessages[1])
		assert.Equal(t, "Квартир стало больше на 2", sentMessages[2])
		assert.Equal(t, "Квартир стало меньше на 3", sentMessages[3])
	})
}
