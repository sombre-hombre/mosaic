package cmd

import (
	"log"

	"github.com/disintegration/imaging"
	"github.com/sombre-hombre/mosaic/app/mosaic"
	"github.com/sombre-hombre/mosaic/app/tiles"
	"github.com/spf13/cobra"
)

// create mosaic:
// mosaic create --in src.jpg --out mosaic.jpg --library tiles/50x50/ --distance redmean

func init() {
	options := struct {
		in, out  string
		library  string
		distance string
	}{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create mosaic from source image using tiles",
		Long: `A photographic mosaic, or a photo-mosaic is a picture (usually a photograph) 
that has been divided into (usually equal sized) rectangular sections, each of which 
is replaced with another picture (called a tile picture).`,
		Run: func(cmd *cobra.Command, args []string) {
			original, err := imaging.Open(options.in)
			if err != nil {
				log.Fatalf("can't open file %s", options.in)
			}

			var dc tiles.DistanceCalculator
			switch options.distance {
			case "redmean":
				dc = tiles.ColorDistanceRedmean
			case "euclidean":
				dc = tiles.ColorDistanceEuclidean
			default:
				log.Fatalf("unknown color distance algorithm '%s'", options.distance)
			}

			lib, err := tiles.NewLibrary(options.library, 50, dc)
			if err != nil {
				log.Fatalf("can't create tile library: %v", err)
			}

			mosaic, err := mosaic.Create(original, *lib)
			if err != nil {
				log.Fatal(err)
			}

			if err = imaging.Save(mosaic, options.out); err != nil {
				log.Fatalf("can't save image: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&options.in, "in", "i", "", "Source image")
	var _ = cmd.MarkFlagRequired("in")
	cmd.Flags().StringVarP(&options.out, "out", "o", "", "Target image")
	_ = cmd.MarkFlagRequired("out")
	cmd.Flags().StringVarP(&options.library, "library", "l", "tiles/50x50/", "Path to tiles library")
	cmd.Flags().StringVarP(&options.distance, "distance", "d", "redmean", "Color distance algorithm: euclidean or redmean")

	rootCmd.AddCommand(cmd)
}
