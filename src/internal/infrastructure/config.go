package infrastructure

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ParseInterval time.Duration
	//ParserStartTgMessage string
}

func NewConfig() *Config {
	parseIntervalMs, err := strconv.ParseInt(os.Getenv("PARSER_INTERVAL_MS"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		ParseInterval: time.Duration(parseIntervalMs) * time.Millisecond,
		//ParserStartTgMessage: os.Getenv("PARSER_START_TG_MESSAGE"),
	}
}
