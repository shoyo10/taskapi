package cmd

import (
	"os"
	"taskapi/pkg/config"
	"taskapi/pkg/zerolog"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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
	log.Debug().Msgf("Config: %+v", cfg)
}
