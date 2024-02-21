package fuze

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	s *http.Server
	*Controller
}

const (
	PORT          = 3000
	READ_TIMEOUT  = time.Second * 2
	WRITE_TIMEOUT = time.Second * 3
)

func NewApp(opts ...Opt) *App {
	p := Params{
		addr:         fmt.Sprintf(":%d", PORT),
		readTimeout:  READ_TIMEOUT,
		writeTimeout: WRITE_TIMEOUT,
	}

	for _, opt := range opts {
		opt(&p)
	}

	c := NewController()
	r := NewRouter(c)

	s := &http.Server{
		Addr:         p.addr,
		ReadTimeout:  p.readTimeout,
		WriteTimeout: p.writeTimeout,
		Handler:      r,
	}

	return &App{
		s:          s,
		Controller: c,
	}
}

func (a *App) Run() error {
	fmt.Printf("App started \nAddr: %s \nHandlers: %d \n", a.s.Addr, a.getHandlersAmount())
	if err := a.s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return nil
}

func (a *App) GracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := a.s.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

type Params struct {
	addr string

	readTimeout  time.Duration
	writeTimeout time.Duration
}

type Opt func(params *Params)

// WithTimeouts sets read/write timeouts for the server
// By default read timeout = 2s write timeout = 3s
func WithTimeouts(read, write time.Duration) Opt {
	return func(p *Params) {
		p.readTimeout = read
		p.writeTimeout = write
	}
}

// WithAddr sets address for the server
func WithAddr(addr string) Opt {
	return func(p *Params) {
		p.addr = addr
	}
}
