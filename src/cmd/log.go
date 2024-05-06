package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogSetting(loglevel string, logformat string) {
	//trace, debug, warn, info, fatal, panic, absent, disable
	switch loglevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "absent":
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case "disable":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	switch logformat {
	case "pretty":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.TimeFieldFormat = time.DateTime
}
