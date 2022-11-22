package cmd

import (
	"context"
	"github.com/sombre-hombre/mosaic/app/server"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	options := struct {
		port int
	}{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start web API service",
		Long:  "TODO: long description",
		Run: func(c *cobra.Command, args []string) {
			// The HTTP Server
			srv := server.NewServer("0.0.0.0", options.port)

			// Server run context
			serverCtx, serverStopCtx := context.WithCancel(context.Background())

			// Listen for syscall signals for process to interrupt/quit
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			go func() {
				<-sig

				// Shutdown signal with grace period of 30 seconds
				shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

				go func() {
					<-shutdownCtx.Done()
					if shutdownCtx.Err() == context.DeadlineExceeded {
						log.Fatal("graceful shutdown timed out.. forcing exit.")
					}
				}()

				// Trigger graceful shutdown
				err := srv.Shutdown(shutdownCtx)
				if err != nil {
					log.Fatal(err)
				}
				serverStopCtx()
			}()

			// Run the server
			err := srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}

			// Wait for server context to be stopped
			<-serverCtx.Done()
		},
	}

	cmd.Flags().IntVarP(&options.port, "port", "p", 8080, "Port")

	rootCmd.AddCommand(cmd)
}
