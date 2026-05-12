package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/roma-glushko/bp/frontend"
	"github.com/roma-glushko/bp/internal/config"
	"github.com/roma-glushko/bp/internal/server"
	"github.com/roma-glushko/bp/internal/storage"
)

func ServeCommand(args []string) error {
	defaults := config.DefaultConfig()

	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", defaults.Port, "port to listen on")
	dataDir := fs.String("data", defaults.DataDir, "data directory path")
	noOpen := fs.Bool("no-open", defaults.NoOpen, "do not open browser automatically")
	debug := fs.Bool("debug", defaults.Debug, "enable debug logging")

	if err := fs.Parse(args); err != nil {
		return err
	}

	level := slog.LevelInfo
	if *debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	store, err := storage.NewTomlStore(*dataDir)
	if err != nil {
		return fmt.Errorf("initializing storage: %w", err)
	}

	mux := server.NewMux(frontend.DistFS, store)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("Blood Pressure Journal is running.")
	fmt.Println()
	fmt.Printf("  Local UI: http://%s\n", addr)
	fmt.Printf("  Data:     %s\n", *dataDir)
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop.")

	slog.Debug("server configured", "port", *port, "data_dir", *dataDir, "no_open", *noOpen)

	if !*noOpen {
		openBrowser(fmt.Sprintf("http://%s", addr))
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-sig
	fmt.Println("\nShutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

func openBrowser(url string) {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	default:
		return
	}

	if err := exec.Command(cmd, url).Start(); err != nil {
		slog.Debug("failed to open browser", "error", err)
	}
}
