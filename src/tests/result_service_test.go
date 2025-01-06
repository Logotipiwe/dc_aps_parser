package tests

import (
	"github.com/stretchr/testify/assert"
	"ports-adapters-study/src/internal/adapters/output"
	"ports-adapters-study/src/internal/core/application"
	"ports-adapters-study/src/internal/core/domain"
	"ports-adapters-study/src/internal/core/ports"
	"reflect"

	"testing"
)

func Test_GetResult(t *testing.T) {

	tests := []struct {
		name    string
		want    *domain.ParseResult
		wantErr bool
	}{
		{name: "kek", want: &domain.ParseResult{ApsNum: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultStorageMock := output.NewResultStorageMock()
			targetClientAdapterMock := output.NewTargetClientAdapterMock([]*domain.ParseResult{tt.want})
			notificationAdapterMock := output.NewNotificationAdapterMock()
			app := application.NewApp(ports.OutputPorts{
				ResultStoragePort: resultStorageMock,
				TargetClientPort:  targetClientAdapterMock,
				NotificationPort:  notificationAdapterMock,
			})
			got, err := app.ResultService.GetResult()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CreateParser(t *testing.T) {
	t.Run("Parser creates", func(t *testing.T) {
		resultStorageMock := output.NewResultStorageMock()
		targetClientAdapterMock := output.NewTargetClientAdapterMock([]*domain.ParseResult{
			{ID: 0, ApsNum: 0},
		})
		notificationAdapterMock := output.NewNotificationAdapterMock()
		app := application.NewApp(ports.OutputPorts{
			ResultStoragePort: resultStorageMock,
			TargetClientPort:  targetClientAdapterMock,
			NotificationPort:  notificationAdapterMock,
		})

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

		app.ParserService.StopAllParsers()
	})

	t.Run("Parser works", func(t *testing.T) {

		resultStorageMock := output.NewResultStorageMock()
		targetClientAdapterMock := output.NewTargetClientAdapterMock([]*domain.ParseResult{
			{ID: 0, ApsNum: 100},
			{ID: 0, ApsNum: 101},
		})
		notificationAdapterMock := output.NewNotificationAdapterMock()
		app := application.NewApp(ports.OutputPorts{
			ResultStoragePort: resultStorageMock,
			TargetClientPort:  targetClientAdapterMock,
			NotificationPort:  notificationAdapterMock,
		})

		notificationAdapterMock.SetCalls(2)

		_, err := app.ParserService.NewParser()

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		notificationAdapterMock.WaitForCalls()

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, 2, len(sentMessages))
		assert.Equal(t, "Квартир стало больше на 1", sentMessages[1])
	})
}
