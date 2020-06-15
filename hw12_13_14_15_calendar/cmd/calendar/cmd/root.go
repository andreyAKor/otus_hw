package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/configs"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/grpc"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/http"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/logging"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository/psql"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	ErrUnknowDatabaseType = errors.New("unknow database type")
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar service application",
	Long:  "The Calendar service is the most simplified service for storing calendar events and sending notifications.",

	RunE: func(cmd *cobra.Command, args []string) error {
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
		var r repository.EventsRepo

		switch c.Database.Type {
		case "memory":
			r = new(memory.Repo)
		case "db":
			rsql := new(psql.Repo)
			if err := rsql.Connect(ctx, c.Database.DSN); err != nil {
				log.Fatal().Err(err).Msg("can't initialize database")
			}
			defer rsql.Close()
			r = rsql
		default:
			log.Fatal().Err(ErrUnknowDatabaseType).Send()
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

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&cfgFile, "config", "", "config file")
	if err := cobra.MarkFlagRequired(pf, "config"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
