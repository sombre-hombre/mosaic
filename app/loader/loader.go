package loader

import (
	"context"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// ImageLoader loads images for tiles library
type ImageLoader interface {
	LoadImages(ctx context.Context, query string, count int, path string) error
}

func getFileExt(u string) string {
	URL, err := url.Parse(u)
	if err != nil {
		return ""
	}

	return path.Ext(URL.Path)
}

func downloadAndSaveImage(ctx context.Context, imgURL, filename string, save imageProcessor) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgURL, http.NoBody)
	if err != nil {
		return err
	}

	client := http.Client{}
	defer client.CloseIdleConnections()

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if ct := res.Header["Content-Type"]; len(ct) > 0 {
		mt, _, err := mime.ParseMediaType(ct[0])
		if err != nil || !strings.HasPrefix(mt, "image/") {
			return errors.New("not an image")
		}
		if path.Ext(filename) == "" {
			exts, err := mime.ExtensionsByType(mt)
			if err == nil {
				filename += exts[0]
			}
		}
	}

	return save(filename, res.Body)
}

type imageProcessor func(filename string, data io.Reader) error

func saveImageToDisk(filename string, data io.Reader) error {
	file, err := os.Create(filename) //nolint:gosec // harmless
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}

	return nil
}
