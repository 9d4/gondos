package cmd

import (
	"context"
	"gondos/api"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type serveCmd struct {
	cmd  *cobra.Command
	opts struct {
		address            string
		corsAllowedOrigins string
	}
}

func newServeCmd() *serveCmd {
	root := &serveCmd{}
	root.cmd = &cobra.Command{
		Use:   "serve",
		Short: "Start http server",
		Run:   root.run,
	}
	root.cmd.Flags().StringVarP(&root.opts.address, "address", "a", "127.0.0.1:8888", "address to listen on")
	return root
}

func (c *serveCmd) run(cmd *cobra.Command, args []string) {
	handler := api.NewHandler()

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
