package main

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/F0rzend/simbirsoft_contest/src/api"
)

const ADDRESS = ":8080"

func main() {
	logger := log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().
		Logger()

	apiServer := api.NewServer()

	server := http.Server{
		Addr:              ":8090",
		Handler:           apiServer.GetHTTPHandler(),
		ReadHeaderTimeout: 3 * time.Second, //nolint:gomnd
	}

	logger.Info().Msgf("Server is running on %q", ADDRESS)

	if err := server.ListenAndServe(); err != nil {
		logger.Error().Err(err).Msg("error while starting server for listening")
	}
}
