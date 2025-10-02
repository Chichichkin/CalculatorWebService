package calculator

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"CalculatorWebService/calculator/storage"
	"CalculatorWebService/internal/config"
	"CalculatorWebService/internal/logger"
	"CalculatorWebService/internal/metrics"
)

type Service struct {
	router  *gin.Engine
	handler *Handler
	metrics *metrics.Metrics
	server  *http.Server
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
	}

	return server
}

func (s *Service) Start() error {
	logger.LogInfo("Server starting", logrus.Fields{
		"address": s.server.Addr,
	})
	return s.server.ListenAndServe()
}

func (s *Service) Shutdown(ctx context.Context) error {
	err := s.handler.Storage.Close()
	if err != nil {
		logger.LogError("Error closing storage", err)
	}
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Service) setupRoutes() {
	s.router.POST("/calculate/addition", s.handler.Addition)
	s.router.POST("/calculate/subtraction", s.handler.Subtraction)
	s.router.POST("calculate/multiplication", s.handler.Multiplication)
	s.router.POST("calculate/division", s.handler.Division)

	s.router.GET("/metrics", gin.WrapH(*s.metrics.Handler))
}
