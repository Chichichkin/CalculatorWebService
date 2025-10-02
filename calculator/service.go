package calculator

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"CalculatorWebService/calculator/storage"
	"CalculatorWebService/internal/config"
	"CalculatorWebService/internal/logger"
	"CalculatorWebService/internal/metrics"
)

// Service struct represents the calculator service with its router, handler, metrics, and HTTP server.
// I prefer to call it service here because we could have multiple services within single server
// Server would be a separate entity that could host multiple services. But for now it's combined within Calculator Service
type Service struct {
	router  *gin.Engine
	handler *Handler
	metrics *metrics.Metrics
	server  *http.Server
	config  config.CalculatorConfig
}

func NewService(configs config.Configs) *Service {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// a bit ugly
	metricsConfig := configs[config.MetricsConfigKey].(config.MetricsConfig)
	serviceConfig := configs[config.CalculatorConfigKey].(config.CalculatorConfig)

	metricsConfig.ServiceName = "calculator" // in a real scenario, this might come from a constant + instance identifier
	metricsConfig.ServiceVersion = serviceConfig.Version

	newMetrics := metrics.NewMetrics(metricsConfig)
	router.Use(logger.LoggingMiddleware())
	router.Use(newMetrics.PrometheusMiddleware())
	router.Use(gin.Recovery())
	newStorage, err := storage.NewStorage(serviceConfig.StorageType, serviceConfig.StorageFilePath)
	if err != nil {
	}
	handler := NewCalculationHandler(newStorage)

	server := &Service{
		router:  router,
		handler: handler,
		metrics: newMetrics,
		server: &http.Server{
			Addr:         ":" + serviceConfig.Port,
			Handler:      router,
			ReadTimeout:  serviceConfig.ReadTimeout,
			WriteTimeout: serviceConfig.WriteTimeout,
			IdleTimeout:  serviceConfig.IdleTimeout,
		},
		config: serviceConfig,
	}
	server.setupRoutes()
	return server
}

func (s *Service) Start() error {
	logger.LogInfo("Calculator starting", logrus.Fields{
		"address": s.server.Addr,
	})
	return s.server.ListenAndServe()
}

func (s *Service) Shutdown(ctx context.Context) {
	logger.LogInfo("Shutting down calculator...")

	logger.LogInfo("Saving records in file...")
	err := s.handler.Storage.Close()
	if err != nil {
		logger.LogError("Error closing storage", err)
	} else {
		logger.LogInfo("Records saved successfully")
	}

	if s.server != nil {
		err := s.server.Shutdown(ctx)
		if err != nil {
			logger.LogError("Calculator forced to shutdown", err)
		} else {
			logger.LogInfo("Calculator shutdown complete")
		}
	}
	logger.LogInfo("Calculator shutdown complete")
}

func (s *Service) setupRoutes() {
	s.router.POST("/calculate/addition", s.handler.Addition)
	s.router.POST("/calculate/subtraction", s.handler.Subtraction)
	s.router.POST("/calculate/multiplication", s.handler.Multiplication)
	s.router.POST("/calculate/division", s.handler.Division)
	s.router.GET("/calculate/recent", s.handler.GetRecentCalculations)

	s.router.GET("/metrics", gin.WrapH(*s.metrics.Handler))
	s.router.GET("/health", s.HealthCheck)
}

func (s *Service) HealthCheck(c *gin.Context) {
	response := struct {
		Status    string    `json:"status"`
		Timestamp time.Time `json:"timestamp"`
		Service   string    `json:"service"`
		Version   string    `json:"version"`
		Uptime    string    `json:"uptime"`
	}{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "calculator",
		Version:   s.config.Version,
	}

	c.JSON(http.StatusOK, response)
}
