// Package serve provides a local development server with live reload.
package serve

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// Server serves the built site locally for development.
type Server struct {
	addr string
	dir  string
}

// New creates a development server for the given output directory.
func New(addr, distDir string) *Server {
	if addr == "" {
		addr = ":1313"
	}
	return &Server{addr: addr, dir: distDir}
}

// Run starts the HTTP server and blocks until interrupted.
func (s *Server) Run() error {
	abs, err := filepath.Abs(s.dir)
	if err != nil {
		return fmt.Errorf("serve: resolve path: %w", err)
	}

	fs := http.FileServer(http.Dir(abs))
	http.Handle("/", fs)

	fmt.Printf("🌐 Serving %s at http://localhost%s\n", abs, s.addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(s.addr, nil); err != nil {
			fmt.Fprintf(os.Stderr, "serve: %v\n", err)
			stop <- syscall.SIGTERM
		}
	}()

	<-stop
	fmt.Println("\n⛔ Server stopped")
	return nil
}
