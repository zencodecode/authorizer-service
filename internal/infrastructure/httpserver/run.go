package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

func Run(ctx context.Context, srv *http.Server, logger service.Logger, name string) error {
	errCh := make(chan error, 1)

	go func() {
		logger.Info(ctx, "server listening", "interface", name, "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("%s server error: %w", name, err)
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Info(ctx, "shutting down server", "interface", name)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
