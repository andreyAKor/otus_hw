package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/configs"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/grpc"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/http"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/logging"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	repo "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/repository"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar service application",
	Long:  "The Calendar service is the most simplified service for storing calendar events and sending notifications.",
	RunE:  run,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&cfgFile, "config", "", "config file")
	if err := cobra.MarkFlagRequired(pf, "config"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init config
	c, err := configs.Init(cfgFile)
	if err != nil {
		return errors.Wrap(err, "init config failed")
	}

	// Init logger
	l := logging.New(c.Logging.File, c.Logging.Level)
	if err := l.Init(); err != nil {
		return errors.Wrap(err, "init logging failed")
	}
	defer l.Close()

	// Init database type
	r, err := repository.New(ctx, c.Database.Type, c.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	if v, ok := r.(repo.DBEventsRepo); ok {
		defer v.Close()
	}

	// Init http-server
	httpSrv, err := http.New(r, c.HTTP.Host, c.HTTP.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize http-server")
	}

	// Init grpc-server
	grpcSrv, err := grpc.New(r, c.GRPC.Host, c.GRPC.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize grpc-server")
	}

	// Init and run app
	a, err := app.New(httpSrv, grpcSrv)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize app")
	}
	if err := a.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("app runnign fail")
	}

	// Graceful shutdown
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)
	<-interruptCh

	if err := a.Close(); err != nil {
		log.Fatal().Err(err).Msg("app closing fail")
	}

	return nil
}
