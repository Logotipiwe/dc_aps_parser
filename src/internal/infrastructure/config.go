package infrastructure

import (
	"dc-aps-parser/src/pkg"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ParseInterval                 time.Duration
	TgAdminChatId                 int64
	TgUserStartMessage            string
	TgParserLaunchMessage         string
	TgParserAlreadyStoppedMessage string
	TgErrorStoppingParserMessage  string
	TgParserStoppedMessage        string
	TgAdminHelpMessage            string
	TgUserHelpMessage             string
	TgStoppedParserStatusMessage  string
	TgUnknownCommandMessage       string
	TgActiveParserStatus          string
	TgInitialApsCountFormat       string
}

func NewConfig() *Config {
	parseIntervalMs, err := strconv.ParseInt(os.Getenv("PARSER_INTERVAL_MS"), 10, 64)
	if err != nil {
		log.Fatal("Error getting parser interval\n", err)
	}
	return &Config{
		ParseInterval:                 time.Duration(parseIntervalMs) * time.Millisecond,
		TgAdminChatId:                 pkg.OsGetInt64NonEmpty("TG_ADMIN_CHAT_ID"),
		TgUserStartMessage:            pkg.OsGetNonEmpty("TG_USER_START_MESSAGE"),
		TgParserLaunchMessage:         pkg.OsGetNonEmpty("TG_PARSER_LAUNCH_MESSAGE"),
		TgParserAlreadyStoppedMessage: pkg.OsGetNonEmpty("TG_PARSER_ALREADY_STOPPED_MESSAGE"),
		TgErrorStoppingParserMessage:  pkg.OsGetNonEmpty("TG_ERROR_STOPPING_PARSER_MESSAGE"),
		TgParserStoppedMessage:        pkg.OsGetNonEmpty("TG_PARSER_STOPPED_MESSAGE"),
		TgAdminHelpMessage:            pkg.OsGetNonEmpty("TG_ADMIN_HELP_MESSAGE"),
		TgUserHelpMessage:             pkg.OsGetNonEmpty("TG_USER_HELP_MESSAGE"),
		TgStoppedParserStatusMessage:  pkg.OsGetNonEmpty("TG_STOPPED_PARSER_STATUS_MESSAGE"),
		TgUnknownCommandMessage:       pkg.OsGetNonEmpty("TG_UNKNOWN_COMMAND_MESSAGE"),
		TgActiveParserStatus:          pkg.OsGetNonEmpty("TG_ACTIVE_PARSER_STATUS_MESSAGE"),
		TgInitialApsCountFormat:       pkg.OsGetNonEmpty("TG_INITIAL_APS_COUNT_FORMAT"),
	}
}
