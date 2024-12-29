package dto

import "ports-adapters-study/src/internal/core/domain"

type ParseResult struct {
	ID     int `json:"id"`
	ApsNum int `json:"aps_num"`
}

func ToResultDto(result domain.ParseResult) ParseResult {
	return ParseResult{
		ID:     result.ID,
		ApsNum: result.ApsNum,
	}
}
