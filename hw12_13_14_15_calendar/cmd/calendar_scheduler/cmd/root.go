package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/calendar"
	configsScheduler "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/configs/scheduler"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/logging"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq/producer/events"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "calendar_scheduler",
	Short: "Calendar scheduler service application",
	Long:  "The Calendar scheduler service is the most simplified service scheduler for searching event to sending notify via RabbitMQ.",
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

//nolint:funlen
func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init config
	c := new(configsScheduler.Config)
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

	// Init events-producer RabbitMQ
	eventsProdMq, err := rmq.New(
		c.RMQ.URI,
		c.Queues.Events.ExchangeName,
		c.Queues.Events.ExchangeType,
		c.Queues.Events.QueueName,
		c.Queues.Events.BindingKey,
		c.RMQ.ReConnect.MaxElapsedTime,
		c.RMQ.ReConnect.InitialInterval,
		c.RMQ.ReConnect.Multiplier,
		c.RMQ.ReConnect.MaxInterval,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize rmq for events-producer")
	}

	// Init events-producer
	eventsProd, err := events.New(
		calendar,
		eventsProdMq,
		c.Producer.CheckEventsToPublishInterval,
		c.Producer.CheckOldEventsInterval,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize events-producer")
	}

	// Init and run app
	a, err := app.New(eventsProd)
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
