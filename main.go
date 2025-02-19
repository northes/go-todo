package main

import (
 "github.com/ichtrojan/go-todo/routes"
 "github.com/joho/godotenv"
 "github.com/rs/zerolog"
 "github.com/rs/zerolog/log"
 "net/http"
 "os"
 "path/filepath"
 "time"
)

func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logLevel := zerolog.InfoLevel
	if os.Getenv("ENV") == "development" {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	logFile, err := os.OpenFile(
		filepath.Join("logs", "error.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}

func main() {
	initLogger()

	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("No .env file found")
		os.Exit(1)
	}

	port, exist := os.LookupEnv("PORT")
	if !exist {
		log.Error().Msg("PORT not set in .env")
		os.Exit(1)
	}

	log.Info().Msgf("Starting server on port %s", port)

	if err := http.ListenAndServe(":"+port, routes.Init()); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}