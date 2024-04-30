package src

import (
	"deployto/src/types"
	"os"
	"time"

	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogSetting(stg types.Settings) {
	//trace, debug, warn, info, fatal, panic, absent, disable
	switch stg.Loglevel {
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

	switch stg.Logformat {
	case "pretty":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	zerolog.TimeFieldFormat = time.DateTime

	pp.Default.SetColoringEnabled(false)
}
