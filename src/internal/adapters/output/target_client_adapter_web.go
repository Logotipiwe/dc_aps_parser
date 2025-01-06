package output

import (
	"ports-adapters-study/src/internal/core/domain"
	krishawebclient "ports-adapters-study/src/internal/infrastructure/krisha"
)

type TargetClientAdapterWeb struct {
	krishaClient *krishawebclient.KrishaWebClient
}

func NewTargetClientWebAdapter() *TargetClientAdapterWeb {
	return &TargetClientAdapterWeb{
		krishawebclient.NewKrishaWebClient(),
	}
}

func (k *TargetClientAdapterWeb) GetParseResult() (*domain.ParseResult, error) {
	mapData, err := k.krishaClient.RequestMapData()
	if err != nil {
		return nil, err
	}
	parseResult := domain.ParseResult{
		ID:     0,
		ApsNum: mapData.NbTotal,
	}
	return &parseResult, err
}