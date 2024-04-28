package src

import (
	"os"

	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogSetting() {
	//TODO set level from flag
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Debug().Msg("zerolog.DebugLevel")
	pp.Default.SetColoringEnabled(false)
}
