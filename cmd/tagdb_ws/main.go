package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "dev.azure.com/trayport/Hackathon/_git/Q/internal/env"
	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
	"dev.azure.com/trayport/Hackathon/_git/Q/internal/tagdb"
)

type config struct {
	portNumber                      int
	storageRoot                     string
	storageWalRollAfterBytes        int64
	storageBackgroundTaskIntervalMs int
}

func main() {
	logger.Info("bootstrapping web server.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := getConfig()

	addSignalHandlers(cancel)
	startDatabase(config, ctx)
	addApiEndpoints()
	addStaticSite()

	if err := runWebServer(config, ctx); err != nil {
		logger.Fatalf("web server exited because %s", err)
	}
}

// Listens for OS signals and initiates shutdown when received.
func addSignalHandlers(cancel context.CancelFunc) {
	logger.Info("adding signal handlers")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		stopSignal := <-stop
		logger.Infof("received OS signal `%s`, shutting down app", stopSignal)
		cancel()
	}()
}

func startDatabase(config config, ctx context.Context) {
	tagdb.Start(
		config.storageRoot,
		ctx,
		tagdb.WithRollAfterBytes(config.storageWalRollAfterBytes),
		tagdb.WithBackgroundTaskIntervalMs(config.storageBackgroundTaskIntervalMs))
}

// Adds handlers for API endpoints.
func addApiEndpoints() {
	logger.Info("adding API endpoint handlers")

	http.HandleFunc("GET /api/keys", List)
	http.HandleFunc("POST /api/keys", Set)
	http.HandleFunc("GET /api/keys/{key}", Get)
	http.HandleFunc("DELETE /api/keys/{key}", Delete)
	http.HandleFunc("POST /api/tags", Tag)
	http.HandleFunc("DELETE /api/tags/{tag}/{key}", Untag)
}

// Adds a handler for static site content.
func addStaticSite() {
	logger.Info("adding static site")
	http.Handle("/", http.FileServer(http.Dir("web")))
}

// Starts the web server.
// Handlers must be added before calling this function.
func runWebServer(config config, ctx context.Context) error {
	port := fmt.Sprintf(":%d", config.portNumber)
	logger.Infof("starting web server on http://localhost%s", port)

	var webErr error
	webServer := &http.Server{Addr: port, Handler: nil}
	go func() {
		if webErr = webServer.ListenAndServe(); webErr != nil {
			return
		}
	}()

	<-ctx.Done()
	if webErr != nil && webErr != http.ErrServerClosed {
		return webErr
	}

	logger.Info("shutting down web server")
	if err := webServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down web server: %s", err)
	}

	return nil
}

// Reads the config from environment variables.
func getConfig() config {
	// Get port.
	portStr := os.Getenv("TAGDB_PORT")
	portNumber, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Warnf("invalid TAGDB_PORT value `%s`", portStr)
		portNumber = 8080
	}

	// WAL roll after bytes.
	walRollAfterBytesStr := os.Getenv("TAGDB_STORAGE_WAL_ROLL_AFTER_BYTES")
	walRollAfterBytes, err := strconv.ParseInt(walRollAfterBytesStr, 10, 64)
	if err != nil {
		logger.Panicf("invalid TAGDB_STORAGE_WAL_ROLL_AFTER_BYTES value `%s`", walRollAfterBytesStr)
	}

	// Background task interval ms.
	backgroundTaskIntervalMsStr := os.Getenv("TAGDB_STORAGE_BACKGROUND_TASK_INTERVAL_MS")
	backgroundTaskIntervalMs, err := strconv.ParseInt(backgroundTaskIntervalMsStr, 10, 64)
	if err != nil {
		logger.Panicf("invalid TAGDB_STORAGE_BACKGROUND_TASK_INTERVAL_MS value `%s`", backgroundTaskIntervalMsStr)
	}

	// Get storage root.
	storageRoot := os.Getenv("TAGDB_STORAGE_ROOT")
	if storageRoot == "" {
		logger.Panicf("cannot start tagDb because TAGDB_STORAGE_ROOT is required")
	}

	info, err := os.Stat(storageRoot)
	if err != nil {
		logger.Panicf("cannot validate TAGDB_STORAGE_ROOT `%s` because %s", storageRoot, err)
	}

	if !info.IsDir() {
		logger.Panicf("cannot start tagDb because TAGDB_STORAGE_ROOT `%s` is not a directory", storageRoot)
	}

	// Success.
	return config{
		portNumber:                      portNumber,
		storageRoot:                     storageRoot,
		storageWalRollAfterBytes:        walRollAfterBytes,
		storageBackgroundTaskIntervalMs: int(backgroundTaskIntervalMs),
	}
}
