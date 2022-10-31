package loader

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadImages(t *testing.T) {
	loader := NewBingLoader("8d81df9b06msh09139a282d05b6ep176e13jsn4f63d0faa31a", 5).(*rapidapiBingLoader)

	expectedCount := 20
	savedFiles := make([]string, 0, expectedCount)
	fakeSaver := func(filename string, data io.Reader) error {
		savedFiles = append(savedFiles, filename)
		return nil
	}

	loader.processor = fakeSaver

	err := loader.LoadImages(context.Background(), "abstract art", expectedCount, "images/")

	require.NoError(t, err)
	require.NotEmpty(t, savedFiles)
	require.LessOrEqual(t, len(savedFiles), expectedCount)

	for _, fn := range savedFiles {
		require.True(t, strings.HasPrefix(fn, "images/"))
	}
}

func Test_getFileExt(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"empty string": {
			input:    "",
			expected: "",
		},
		"no extension": {
			input:    "https://ya.ru/icon",
			expected: "",
		},
		"no file": {
			input:    "https://ya.ru/",
			expected: "",
		},
		"no path": {
			input:    "https://ya.ru",
			expected: "",
		},
		"extension": {
			input:    "https://ya.ru/icon.png",
			expected: ".png",
		},
		"invalid url": {
			input:    "not an url",
			expected: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := getFileExt(tc.input)

			require.Equal(t, tc.expected, actual)
		})
	}
}
