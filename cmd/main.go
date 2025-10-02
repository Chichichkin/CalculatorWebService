package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"CalculatorWebService/calculator"
	"CalculatorWebService/internal/config"
	"CalculatorWebService/internal/logger"
)

func main() {
	requiredConfigs := []string{
		config.MetricsConfigKey,
		config.LoggerConfigKey,
		config.CalculatorConfigKey,
	}
	configs := config.LoadConfigs(requiredConfigs)

	logger.InitLogger(configs[config.LoggerConfigKey].(config.LoggerConfig))

	srv := calculator.NewService(configs)

	go GracefulShutdown(srv)

	if err := srv.Start(); err != nil {
		logger.LogError("Server error", err)
	}

}

func GracefulShutdown(srv *calculator.Service) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.LogInfo("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.Shutdown(ctx)

	logger.LogInfo("Server exited")
}
