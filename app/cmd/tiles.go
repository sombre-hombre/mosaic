package cmd

import (
	"log"

	"github.com/sombre-hombre/mosaic/app/tiles"
	"github.com/spf13/cobra"
)

// mosaic tiles --in images/cats --out tiles/cats --size 50

func init() {
	opts := struct {
		in, out  string
		tileSize int
	}{}

	cmd := &cobra.Command{
		Use:   "tiles",
		Short: "Make tiles from arbitrary images",
		Long:  "Make required sized square tiles from arbitrary images located in specified folder.",
		Run: func(cmd *cobra.Command, args []string) {
			err := tiles.PrepareTiles(opts.in, opts.out, opts.tileSize)
			if err != nil {
				log.Fatalf("can't create tiles: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.in, "in", "i", "", "folder with source images")
	cmd.MarkFlagRequired("in")
	cmd.Flags().StringVarP(&opts.out, "out", "o", "", "output folder for tiles")
	cmd.MarkFlagRequired("out")
	cmd.Flags().IntVarP(&opts.tileSize, "size", "s", 50, "tile size (px)")

	rootCmd.AddCommand(cmd)
}
