package tiles

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"

	"github.com/disintegration/imaging"
)

const (
	// DefaultTileSize — default tile size px
	DefaultTileSize = 50
)

type Library struct {
	// Directory with prepared tiles
	tilesDir string
	tiles    map[string]AvgColor
	tileSize int
	distance DistanceCalculator
}

func NewLibrary(tilesDir string, tileSize int, dc DistanceCalculator) (*Library, error) {
	l := &Library{
		tilesDir: tilesDir,
		tileSize: tileSize,
		distance: dc,
	}

	files, err := ioutil.ReadDir(tilesDir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in %s", tilesDir)
	}

	l.tiles = make(map[string]AvgColor, len(files))
	for _, fi := range files {
		fn := fi.Name()
		img, err := imaging.Open(path.Join(tilesDir, fn))
		if err != nil {
			log.Printf("can't open image %s: %v", fn, err)
			continue
		}

		if img.Bounds().Dx() != l.tileSize || img.Bounds().Dy() != l.tileSize {
			log.Printf("wrong tile size in %s, skip", fn)
			continue
		}

		l.tiles[fn] = NewAvgColor(img)
	}

	return l, nil
}

// TileSize returns size of tiles in library
func (l Library) TileSize() int {
	return l.tileSize
}

func (l Library) FindTile(subImage image.Image) (image.Image, error) {
	avg := NewAvgColor(subImage)
	smallest := math.MaxFloat64
	var tilePath string
	for path, clr := range l.tiles {
		dist := l.distance(avg, clr)
		if dist < smallest {
			smallest = dist
			tilePath = path
		}
	}

	if tilePath == "" {
		return nil, errors.New("tile not found")
	}

	return imaging.Open(path.Join(l.tilesDir, tilePath))
}

// PrepareTiles prepares tiles from images found in sourceDir and stores them to targetDir.
// If tileSize is 0, then DefaultTileSize is used instead.
func PrepareTiles(sourceDir string, targetDir string, tileSize int) error {
	if tileSize == 0 {
		tileSize = DefaultTileSize
	}

	if err := os.MkdirAll(targetDir, 0o700); err != nil { // If path is already a directory, MkdirAll does nothing
		return fmt.Errorf("can't make directory %s: %w", targetDir, err)
	}

	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		img, err := imaging.Open(path.Join(sourceDir, f.Name()))
		if err != nil {
			log.Printf("can't open image %s: %v", f.Name(), err)
			continue
		}

		if img.Bounds().Dx() < tileSize || img.Bounds().Dy() < tileSize {
			// image too small, skip it
			continue
		}

		// resize image preserving aspect ratio
		if img.Bounds().Dx() < img.Bounds().Dy() {
			img = imaging.Resize(img, tileSize, 0, imaging.Lanczos)
		} else {
			img = imaging.Resize(img, 0, tileSize, imaging.Lanczos)
		}

		// crop image to square
		img = imaging.CropAnchor(img, tileSize, tileSize, imaging.TopLeft)

		// save prepared tile to file
		err = imaging.Save(img, path.Join(targetDir, f.Name()))
		if err != nil {
			log.Printf("can't save %s: %v", f.Name(), err)
			continue
		}
	}

	return nil
}
