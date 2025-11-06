package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
	"dev.azure.com/trayport/Hackathon/_git/Q/internal/tagdb"
)

func main() {
	logger.Info("bootstrapping web server.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addSignalHandlers(cancel)
	startDatabase(ctx)
	addApiEndpoints()
	addStaticSite()

	if err := runWebServer(ctx); err != nil {
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

func startDatabase(ctx context.Context) {
	// TODO: Plumbing.
	tagdb.Start("C:\\Users\\davidr\\.q", ctx)
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
func runWebServer(ctx context.Context) error {
	logger.Info("starting web server on http://localhost:31979")

	var webErr error
	webServer := &http.Server{Addr: ":31979", Handler: nil}
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
