package tiles

import (
	"fmt"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/require"
)

func Test_NewAvgColor(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected AvgColor
	}{
		"red": {
			input:    "red.png",
			expected: AvgColor{0xffff, 0, 0},
		},
		"green": {
			input:    "green.png",
			expected: AvgColor{0, 0xffff, 0},
		},
		"blue": {
			input:    "blue.png",
			expected: AvgColor{0, 0, 0xffff},
		},
		"white": {
			input:    "white.png",
			expected: AvgColor{0xffff, 0xffff, 0xffff},
		},
		"black": {
			input:    "black.png",
			expected: AvgColor{0, 0, 0},
		},
		"black & white": {
			input:    "bw.png",
			expected: AvgColor{0xffff / 2.0, 0xffff / 2.0, 0xffff / 2.0},
		},
		"white & black": {
			input:    "wb.png",
			expected: AvgColor{0xffff / 2.0, 0xffff / 2.0, 0xffff / 2.0},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			img, err := imaging.Open(fmt.Sprintf("testdata/%s", tc.input))
			require.NoError(t, err)

			actual := NewAvgColor(img)

			require.Equal(t, tc.expected, actual)
		})
	}
}

func Test_euclidianDistance(t *testing.T) {
	tests := map[string]struct {
		img1, img2 AvgColor
		expected   float64
	}{
		"same image": {
			img1:     AvgColor{0, 0, 0},
			img2:     AvgColor{0, 0, 0},
			expected: 0.0,
		},
		"same image 2": {
			img1:     AvgColor{0xffff / 2.0, 0xffff / 2.0, 0xffff / 2.0},
			img2:     AvgColor{0xffff / 2.0, 0xffff / 2.0, 0xffff / 2.0},
			expected: 0.0,
		},
		"red v blue": {
			img1:     AvgColor{0xffff, 0, 0},
			img2:     AvgColor{0, 0, 0xffff},
			expected: 92680,
		},
		"red v green": {
			img1:     AvgColor{0xffff, 0, 0},
			img2:     AvgColor{0, 0xffff, 0},
			expected: 92680,
		},
		"red v white": {
			img1:     AvgColor{0xffff, 0, 0},
			img2:     AvgColor{0xffff, 0xffff, 0xffff},
			expected: 92680,
		},
		"blue v white": {
			img1:     AvgColor{0, 0, 0xffff},
			img2:     AvgColor{0xffff, 0xffff, 0xffff},
			expected: 92680,
		},
		"green v white": {
			img1:     AvgColor{0, 0xffff, 0},
			img2:     AvgColor{0xffff, 0xffff, 0xffff},
			expected: 92680,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := euclidianDistance(tc.img1, tc.img2)

			require.InDelta(t, tc.expected, actual, .5)
		})
	}
}

func Test_redmeanDistance(t *testing.T) {
	tests := map[string]struct {
		img1, img2 AvgColor
		expected   float64
	}{
		"same image": {
			img1:     AvgColor{0, 0, 0},
			img2:     AvgColor{0, 0, 0},
			expected: 0.0,
		},
		"same image 2": {
			img1:     AvgColor{0xffff / 2.0, 0xffff / 2.0, 0xffff / 2.0},
			img2:     AvgColor{0xffff / 2.0, 0xffff / 2.0, 0xffff / 2.0},
			expected: 0.0,
		},
		"red v blue": {
			img1:     AvgColor{0xffff, 0, 0},
			img2:     AvgColor{0, 0, 0xffff},
			expected: 146540,
		},
		"red v green": {
			img1:     AvgColor{0xffff, 0, 0},
			img2:     AvgColor{0, 0xffff, 0},
			expected: 167082,
		},
		"red v white": {
			img1:     AvgColor{0xffff, 0, 0},
			img2:     AvgColor{0xffff, 0xffff, 0xffff},
			expected: 160527,
		},
		"blue v white": {
			img1:     AvgColor{0, 0, 0xffff},
			img2:     AvgColor{0xffff, 0xffff, 0xffff},
			expected: 167082,
		},
		"green v white": {
			img1:     AvgColor{0, 0xffff, 0},
			img2:     AvgColor{0xffff, 0xffff, 0xffff},
			expected: 146540,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := redMeanDistance(tc.img1, tc.img2)

			require.InDelta(t, tc.expected, actual, .5)
		})
	}
}
