package tiles

import (
	"image"
	"math"
)

// AvgColor is average color of image
type AvgColor [3]float64

// R is red component
func (c AvgColor) R() float64 {
	return c[0]
}

// G is green component
func (c AvgColor) G() float64 {
	return c[1]
}

// B is blue component
func (c AvgColor) B() float64 {
	return c[2]
}

// NewAvgColor returns average color of the img
func NewAvgColor(img image.Image) AvgColor {
	bounds := img.Bounds()
	r, g, b := 0.0, 0.0, 0.0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
		}
	}
	total := float64(bounds.Dx() * bounds.Dy())
	return AvgColor{r / total, g / total, b / total}
}

// Color distance calculator
type DistanceCalculator func(AvgColor, AvgColor) float64

var (
	ColorDistanceEuclidean DistanceCalculator = euclideanDistance
	ColorDistanceRedmean   DistanceCalculator = redMeanDistance
)

// euclideanDistance implements euclidean color distance
func euclideanDistance(c1, c2 AvgColor) float64 {
	dr, dg, db := c1.R()-c2.R(), c1.G()-c2.G(), c1.B()-c2.B()
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// redmeanDistance implements so called "redmean" approximation https://en.wikipedia.org/wiki/Color_difference#sRGB
func redMeanDistance(c1, c2 AvgColor) float64 {
	r := (c1.R() + c2.R()) / 2.0
	dr, dg, db := c1.R()-c2.R(), c1.G()-c2.G(), c1.B()-c2.B()

	return math.Sqrt((2+r/0x10000)*dr*dr + 4*dg*dg + (2+(0xffff-r)/0x10000)*db*db)
}
