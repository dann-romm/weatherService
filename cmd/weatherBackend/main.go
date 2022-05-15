package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"testbackend/internal/appctl"
	"testbackend/internal/weather"
	"time"
)

func logError(err error) {
	if err = fmt.Errorf("%w", err); err != nil {
		fmt.Println(err)
	}
}

type server struct {
	weatherService weather.OpenWeatherMapService
}

func (s *server) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	curWeather, err := s.weatherService.GetWeather()
	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("418 I'm a teapot"))
		return
	}
	_, err = curWeather.WriteTo(w)
	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("418 I'm a teapot"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *server) appStart(ctx context.Context, halt <-chan struct{}) error {
	r := mux.NewRouter()
	r.Handle("/current-weather", s)
	var httpServer = http.Server{
		Addr:              ":8900",
		Handler:           r,
		ReadTimeout:       time.Millisecond * 250,
		ReadHeaderTimeout: time.Millisecond * 200,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Minute * 30,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	var errShutdown = make(chan error, 1)
	go func() {
		defer close(errShutdown)
		select {
		case <-halt:
		case <-ctx.Done():
		}
		if err := httpServer.Shutdown(ctx); err != nil {
			errShutdown <- err
		}
	}()
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	err, ok := <-errShutdown
	if ok {
		return err
	}
	return nil
}

func main() {
	var srv server
	var svc = appctl.ServiceKeeper{
		Services: []appctl.Service{
			&srv.weatherService,
		},
		ShutdownTimeout: time.Second * 10,
		PingPeriod:      time.Millisecond * 500,
	}
	var app = appctl.Application{
		MainFunc:           srv.appStart,
		Resources:          &svc,
		TerminationTimeout: time.Second * 10,
	}
	if err := app.Run(); err != nil {
		logError(err)
		os.Exit(1)
	}
}
