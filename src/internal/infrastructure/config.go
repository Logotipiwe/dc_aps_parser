package infrastructure

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ParseInterval         time.Duration
	TgParserLaunchMessage string
	TgUserStartMessage    string
}

func NewConfig() *Config {
	parseIntervalMs, err := strconv.ParseInt(os.Getenv("PARSER_INTERVAL_MS"), 10, 64)
	if err != nil {
		log.Fatal("Error getting parser interval\n", err)
	}
	return &Config{
		ParseInterval:         time.Duration(parseIntervalMs) * time.Millisecond,
		TgParserLaunchMessage: os.Getenv("TG_PARSER_LAUNCH_MESSAGE"),
		TgUserStartMessage:    os.Getenv("TG_USER_START_MESSAGE"),
	}
}
