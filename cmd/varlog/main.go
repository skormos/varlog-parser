package main

func main() {
	mainLogger := stdoutLoggerContext("main").Logger()

	mainLogger.Info().Msg("started")

	mainLogger.Info().Msg("good-bye")
}
