package cmd

import (
	"context"
	"os"
	"taskapi/internal/delivery/http"
	"taskapi/internal/repository"
	"taskapi/internal/service"
	"taskapi/pkg/config"
	"taskapi/pkg/echorouter"
	"taskapi/pkg/zerolog"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// ServerCmd ...
var ServerCmd = &cobra.Command{
	Run: run,
	Use: "taskapi",
}

func run(cmd *cobra.Command, args []string) {
	cfg, err := config.New()
	if err != nil {
		os.Exit(1)
		return
	}

	zerolog.Init(cfg.Log)

	app := fx.New(
		fx.Supply(*cfg),
		fx.Provide(
			repository.New,
			service.New,
			echorouter.FxNewEcho,
			http.NewHandler,
		),
		fx.Invoke(
			http.SetRoutes,
		),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Error().Msg(err.Error())
		return
	}

	<-app.Done()

	log.Info().Msgf("main: shutting down %s...", cmd.Name())

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Error().Msg(err.Error())
	}
}
