package logger

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func NewLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	debug := flag.Bool("debug", false, "Enable debug mode with more verbose logging")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i any) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i any) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i any) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i any) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log := zerolog.New(output).With().Timestamp().Logger()
	return log
}