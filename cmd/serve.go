package cmd

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"gondos/api"
	"gondos/db"
	"gondos/internal/app"
	store "gondos/stores/mysql"
)

type serveCmd struct {
	cmd  *cobra.Command
	opts struct {
		address            string
		corsAllowedOrigins string
		dsn                string
		dbLogFile          string
		dbLogStdout        bool
	}
}

func newServeCmd() *serveCmd {
	root := &serveCmd{}
	root.cmd = &cobra.Command{
		Use:   "serve",
		Short: "Start http server",
		Run:   root.run,
	}
	root.cmd.Flags().StringVarP(&root.opts.address, "address", "a", env("ADDRESS", "127.0.0.1:8888"), "address to listen on")
	root.cmd.Flags().StringVar(&root.opts.dsn, "dsn", env("DSN"), "database dsn")
	root.cmd.Flags().StringVar(&root.opts.dbLogFile, "db-log", "db.log", "database log file")
	root.cmd.Flags().BoolVar(&root.opts.dbLogStdout, "db-log-stdout", false, "log database to stdout instead")
	return root
}

func (c *serveCmd) run(cmd *cobra.Command, args []string) {
	// setup db
	log.Info().Msg("connecting to database...")
	db_, err := db.Open(c.opts.dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect database")
	}

	db.SetQueryLogger(os.Stdout)
	if !c.opts.dbLogStdout {
		dbLogWriter, err := db.NewLogWriter(c.opts.dbLogFile)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot create db log file")
		}
		db.SetQueryLogger(dbLogWriter)
	}

	schemaName, err := db.ParseSchema(c.opts.dsn)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	app := app.New(app.Dependencies{
		UserStore: store.NewUserStore(db_, schemaName),
		ListStore: store.NewListStore(db_, schemaName),
	})

	handler := api.NewHandler(app)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         c.opts.address,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handler,
	}
	srvErr := make(chan error, 1)

	go func() {
		log.Info().Msgf("Listening on %s", c.opts.address)
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err := <-srvErr:
		log.Fatal().Err(err).Send()
		return
	case <-ctx.Done():
		log.Info().Msg("Stopping")
		stop()
	}

	srv.Shutdown(context.Background())
}
