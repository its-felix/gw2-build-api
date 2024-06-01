package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func listenAddr() string {
	port, _ := strconv.Atoi(os.Getenv("AWS_LWA_PORT"))
	port = cmp.Or(port, 8080)

	return fmt.Sprintf(":%d", port)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	srv := &http.Server{
		Addr:    listenAddr(),
		Handler: setupMux(),
	}

	if err := run(ctx, srv); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, srv *http.Server) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			slog.Error("error shutting down http server", slog.String("err", err.Error()))
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}
