package application

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	return &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Printf("failed to close redis connection: %v", err)
		}
	}()

	fmt.Println("Server is running on port 3000")

	ch := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to listen to server: %w", err)
		}
		close(ch)
	}()

	select {
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		err := server.Shutdown(timeout)
		if err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	case err := <-ch:
		return err
	}
	return nil
}
