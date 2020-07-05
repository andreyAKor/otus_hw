package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/calendar"
	configsCalendar "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/configs/calendar"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/grpc"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/http"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/logging"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar API service application",
	Long:  "The Calendar API service is the most simplified service for storing calendar events and sending notifications.",
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
	c := new(configsCalendar.Config)
	if err := c.Init(cfgFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	// Init logger
	l := logging.New(c.Logging.File, c.Logging.Level)
	if err := l.Init(); err != nil {
		return errors.Wrap(err, "init logging failed")
	}
	defer l.Close()

	// Init calendar
	calendar, err := calendar.New(ctx, c.Database.Type, c.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer calendar.Close()

	// Init http-server
	httpSrv, err := http.New(calendar, c.HTTP.Host, c.HTTP.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize http-server")
	}

	// Init grpc-server
	grpcSrv, err := grpc.New(calendar, c.GRPC.Host, c.GRPC.Port)
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

	log.Info().Msg("Stopping...")

	if err := a.Close(); err != nil {
		log.Fatal().Err(err).Msg("app closing fail")
	}

	log.Info().Msg("Stopped")

	return nil
}
