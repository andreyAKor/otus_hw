package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/app/sender"
	configsSender "github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/configs/sender"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/logging"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/rmq/consumer"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "calendar_sender",
	Short: "Calendar sender service application",
	Long:  "The Calendar sender service is the most simplified service sender for searching event to sending notify via RabbitMQ.",
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
	c := new(configsSender.Config)
	if err := c.Init(cfgFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	// Init logger
	l := logging.New(c.Logging.File, c.Logging.Level)
	if err := l.Init(); err != nil {
		return errors.Wrap(err, "init logging failed")
	}
	defer l.Close()

	// Init RabbitMQ
	mq, err := rmq.New(
		c.RMQ.URI,
		c.RMQ.ExchangeName,
		c.RMQ.ExchangeType,
		c.RMQ.QueueName,
		c.RMQ.BindingKey,
		c.RMQ.ReConnect.MaxElapsedTime,
		c.RMQ.ReConnect.InitialInterval,
		c.RMQ.ReConnect.Multiplier,
		c.RMQ.ReConnect.MaxInterval,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize rmq")
	}

	// Init consumer
	cons, err := consumer.New(
		mq,
		c.Consumer.ConsumerTag,
		c.Consumer.QosPrefetchCount,
		c.Consumer.Threads,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize rmq-consumer")
	}

	// Init and run app
	a, err := app.New(cons)
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
