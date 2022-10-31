package cmd

import (
	"context"
	"log"

	"github.com/sombre-hombre/mosaic/app/loader"
	"github.com/spf13/cobra"
)

func init() {
	options := struct {
		apiKey      string
		concurrency int
		query       string
		out         string
		count       int
	}{}

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load images for tiles from API",
		Long:  "Load images for tiles library from external API, i.e. bing image search",
		Run: func(cmd *cobra.Command, args []string) {
			loader := loader.NewBingLoader(options.apiKey, options.concurrency)
			err := loader.LoadImages(context.Background(), options.query, options.count, options.out)
			if err != nil {
				log.Fatalf("can't load images: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&options.apiKey, "key", "k", "", "Rapidapi api key")
	cmd.MarkFlagRequired("key")
	cmd.Flags().StringVarP(&options.query, "query", "q", "", "Image search query")
	cmd.MarkFlagRequired("query")
	cmd.Flags().StringVarP(&options.out, "out", "o", "images", "Output path")
	cmd.Flags().IntVarP(&options.count, "count", "c", 100, "Images count")
	cmd.Flags().IntVarP(&options.concurrency, "concurrency", "l", 5, "Concurrency level")

	rootCmd.AddCommand(cmd)
}
