package app

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type App struct {
	r    repository.EventsRepo
	host string
	port int
}

func New(r repository.EventsRepo, host string, port int) (*App, error) {
	return &App{r, host, port}, nil
}

func (a *App) Run(ctx context.Context) error {
	// Running application http-server.
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.Hello)

	// middleware
	handler := a.Logger(mux)

	server := &http.Server{
		Addr:    net.JoinHostPort(a.host, strconv.Itoa(a.port)),
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil {
		return errors.Wrap(err, "http-server listen fail")
	}

	return nil
}

// Logger output log info of request, e.g.: r.Method, r.URL etc.
func (a *App) Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newAppResponseWriter(w)

		start := time.Now()
		defer func() {
			i := log.Info()
			i.Str("IP", strings.Split(r.RemoteAddr, ":")[0]).
				Str("StartAt", start.String()).
				Str("Method", r.Method)

			if r.URL != nil {
				i.Str("Path", r.URL.Path)
			}

			i.Str("Proto", r.Proto).
				Int("Status", rw.statusCode).
				TimeDiff("Latency", time.Now(), start)

			if len(r.UserAgent()) > 0 {
				i.Str("UserAgent", r.UserAgent())
			}

			i.Msg("new request")
		}()

		handler.ServeHTTP(rw, r)
	})
}

// The "hello-world" handler.
func (a *App) Hello(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("hello - world")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// App wrapper over http.ResponseWriter.
type appResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newAppResponseWriter(w http.ResponseWriter) *appResponseWriter {
	return &appResponseWriter{w, http.StatusOK}
}

func (a *appResponseWriter) WriteHeader(code int) {
	a.statusCode = code
	a.ResponseWriter.WriteHeader(code)
}
