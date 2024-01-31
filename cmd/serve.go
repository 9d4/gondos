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
	"gondos/app"
	"gondos/db"
)

type serveCmd struct {
	cmd  *cobra.Command
	opts struct {
		address            string
		corsAllowedOrigins string
		dsn                string
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
	return root
}

func (c *serveCmd) run(cmd *cobra.Command, args []string) {
	// setup db
	log.Info().Msg("connecting to database...")
	db_, err := db.Open(c.opts.dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect database")
	}

	dbLogWriter, err := db.NewLogWriter("apk.log")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create db log file")
	}
	db.SetQueryLogger(dbLogWriter)

	app := app.New(&app.Config{
		DB: db_,
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
