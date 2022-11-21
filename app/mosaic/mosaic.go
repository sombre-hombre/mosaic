package mosaic

import (
	"errors"
	"fmt"
	"image"
	"image/draw"

	"github.com/sombre-hombre/mosaic/app/tiles"
)

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

// Create creates mosaic image from original image using title library lib
func Create(original image.Image, lib tiles.Library) (image.Image, error) {
	// image.Image interface does not implement SubImage, but some implementations does,
	// so we try to cast original to subImager to get the method.
	subimgr, ok := original.(subImager)
	if !ok {
		return nil, errors.New("unsupported image type")
	}

	tileSize := lib.TileSize()

	mosaic := image.NewNRGBA(original.Bounds())
	for x := original.Bounds().Min.X; x < original.Bounds().Max.X; x += tileSize {
		for y := original.Bounds().Min.Y; y < original.Bounds().Max.Y; y += tileSize {
			tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
			subimage := subimgr.SubImage(tileBounds)
			tile, err := lib.FindTile(subimage)
			if err != nil {
				return nil, fmt.Errorf("can't find tile: %v", err)
			}

			draw.Draw(mosaic, tileBounds, tile, image.Point{X: 0, Y: 0}, draw.Src)
		}
	}

	return mosaic, nil
}
