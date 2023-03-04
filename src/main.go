package main

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/F0rzend/simbirsoft_contest/src/api"
	"github.com/F0rzend/simbirsoft_contest/src/common"
)

const ADDRESS = ":8080"

func main() {
	config := common.ConfigFromEnv()

	logger := log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().
		Logger()

	apiServer, err := api.NewServer(config)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	handler, err := apiServer.GetHTTPHandler(&logger)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	server := http.Server{
		Addr:              ADDRESS,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Info().Msgf("Server is running on %q", ADDRESS)

	if err := server.ListenAndServe(); err != nil {
		logger.Error().Err(err).Msg("error while starting server for listening")
	}
}
