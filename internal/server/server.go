package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/sunr3d/basic-marketplace/internal/config"
)

type Server struct {
	httpServer *http.Server
	log        *zap.Logger
}

func New(router *gin.Engine, cfg *config.Config, log *zap.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.HTTPHost + ":" + cfg.HTTPPort,
			Handler:      router,
			ReadTimeout:  time.Duration(cfg.HTTPReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.HTTPWriteTimeout) * time.Second,
		},
		log: log,
	}
}

func (s *Server) Start() error {
	go func() {
		s.log.Info("Запуск HTTP сервера", zap.String("addr", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	s.log.Info("Получен сигнал завершения, идет остановка сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Error("Ошибка при остановке сервера", zap.Error(err))
		return err
	}
	s.log.Info("Сервер остановлен успешно")
	return nil
}
