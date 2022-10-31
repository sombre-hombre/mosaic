package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

// NewBingLoader returns rapidapi bing image search loader
func NewBingLoader(apiKey string, concurrency int) ImageLoader {
	if concurrency < 1 {
		concurrency = 5
	}

	client := http.Client{}

	return &rapidapiBingLoader{
		apiKey:              apiKey,
		concurrentDownloads: concurrency,
		processor:           saveImageToDisk,
		client:              client,
	}
}

type rapidapiBingLoader struct {
	apiKey              string
	concurrentDownloads int
	processor           imageProcessor
	client              http.Client
	imgCounter          int
}

// LoadImages loads images from https://rapidapi.com/microsoft-azure-org-microsoft-cognitive-services/api/bing-image-search1/
//
// query — query for image search
//
// count — count of images to load
func (l *rapidapiBingLoader) LoadImages(ctx context.Context, query string, count int, outPath string) error {
	pageSize := 50
	if count < pageSize {
		pageSize = count
	}

	if err := os.MkdirAll(outPath, 0o700); err != nil { // If path is already a directory, MkdirAll does nothing
		return fmt.Errorf("can't make directory %s: %w", outPath, err)
	}

	l.imgCounter = 0
	for offset := 0; offset < count; offset += pageSize {
		searchURL := fmt.Sprintf(
			"https://bing-image-search1.p.rapidapi.com/images/search?q=%s&offset=%d&count=%d",
			url.QueryEscape(query),
			offset,
			pageSize,
		)

		responseData, err := l.getSearchPage(ctx, searchURL)
		if err != nil {
			return err
		}

		l.downloadImages(ctx, responseData, outPath)
	}

	return nil
}

func (l *rapidapiBingLoader) getSearchPage(ctx context.Context, searchURL string) (*imageSearchResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", l.apiKey)
	req.Header.Add("X-RapidAPI-Host", "bing-image-search1.p.rapidapi.com")

	res, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseData := &imageSearchResponse{}
	if err := json.NewDecoder(res.Body).Decode(responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}

// downloadImages concurrently downloads found images to outPath
func (l *rapidapiBingLoader) downloadImages(ctx context.Context, searchResult *imageSearchResponse, outPath string) {
	// TODO: make semaphore package
	wg := sync.WaitGroup{}
	defer wg.Wait()
	sem := make(chan any, l.concurrentDownloads)
	for _, img := range searchResult.Value {
		sem <- struct{}{}
		wg.Add(1)
		imgURL := img.URL
		ext := getFileExt(img.URL)
		filename := path.Join(outPath, fmt.Sprintf("%04d%s", l.imgCounter, ext))
		l.imgCounter++
		go func(wg *sync.WaitGroup) {
			defer func() {
				<-sem
				wg.Done()
			}()
			err := downloadAndSaveImage(ctx, imgURL, filename, l.processor)
			if err != nil {
				log.Printf("can't process image %s: %v", filename, err)
			}
		}(&wg)
	}
}

type imageSearchResponse struct {
	TotalCount int         `json:"totalEstimatedMatches"`
	Value      []imageData `json:"value"`
}

type imageData struct {
	URL    string `json:"contentUrl"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
