package echorouter

import (
	"context"
	"net/http"
	"taskapi/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// Config setting http config
type Config struct {
	Debug        bool   `yaml:"debug"`
	Address      string `yaml:"address"`
	DisablePprof bool   `yaml:"disable_pprof"`
}

// NewEcho create new engine for handler to register
func NewEcho(cfg *Config) *echo.Echo {
	echo.NotFoundHandler = NotFoundHandler
	echo.MethodNotAllowedHandler = NotFoundHandler

	e := echo.New()

	if cfg.Debug {
		e.Debug = true
		e.HideBanner = false
		e.HidePort = false
	} else {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
	}
	e.HTTPErrorHandler = ErrorHandler

	e.Use(MiddlewareRequestID())
	e.Use(MiddlewareAccessLog())
	e.Use(MiddlewareLogWithRequestID())
	e.Use(MiddlewareCorsConfig)
	e.Use(MiddlewareRecover())

	setDefaultRoute(e, cfg)

	return e
}

// NotFoundHandler responds not found response.
func NotFoundHandler(c echo.Context) error {
	return c.String(http.StatusNotFound, "page not found")
}

// ErrorHandler responds error response according to given error.
func ErrorHandler(err error, c echo.Context) {
	echoErr, ok := err.(*echo.HTTPError)
	if ok {
		err = c.JSON(echoErr.Code, echoErr)
		if err != nil {
			log.Err(err).Msgf("%v", err)
		}
		return
	}

	causeErr := errors.Cause(err)
	httpErr := errors.GetHTTPError(causeErr)
	err = c.JSON(httpErr.HTTPCode, httpErr)
	if err != nil {
		log.Err(err).Msgf("%v", err)
	}
}

func setDefaultRoute(e *echo.Echo, cfg *Config) {
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	if !cfg.DisablePprof {
		RegisterPprofRouter(e)
	}
}

// FxNewEcho create new echo server with fx lifecycle
func FxNewEcho(cfg *Config, lc fx.Lifecycle) *echo.Echo {
	e := NewEcho(cfg)
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info().Msgf("Starting echo server, listen on %s", cfg.Address)
			go func() {
				err := e.Start(cfg.Address)
				if err != nil {
					if err == http.ErrServerClosed {
						log.Info().Msg("Echo server closed.")
					} else {
						log.Error().Msgf("Error echo server, err: %s", err.Error())
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Stopping echo server.")
			return e.Shutdown(ctx)
		},
	})
	return e
}
