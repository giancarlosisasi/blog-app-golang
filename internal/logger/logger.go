package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// log.Debug().Msg("This message appears only when log level set to Debug")
	// log.Info().Msg("This message appears when log level set to Debug or Info")
	log.Info().Msg("> ✅ log setup complete ")

}
