package processors

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"

	// Importing image/jpeg and image/png to help decode
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// BlockSize is the abstraction of a widthxHeight map
type BlockSize struct {
	Width  int
	Height int
}

// Rect used for image bounds
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ImageInfo reflects information required by mediafaker to recreate the file. This amounts to width height and pixel info
type ImageInfo struct {
	Width       int      `json:"W"`
	Height      int      `json:"H"`
	PixelInfo   []string `json:"P"`
	BlockWidth  int      `json:"BW"`
	BlockHeight int      `json:"BH"`
}

// PixelRectangle reflects the size, position and color of a pixel rectangle
type PixelRectangle struct {
	Color     string
	Rectangle image.Rectangle
}

// ImageProcessor structure
type ImageProcessor struct {
}

// Inspect Command
func (processor *ImageProcessor) Inspect(sourcePath string) (ImageInfo, error) {
	imageInfo := ImageInfo{}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return imageInfo, err
	}
	img, _, err := image.Decode(sourceFile)
	if err != nil {
		return imageInfo, err
	}
	err = sourceFile.Close()
	if err != nil {
		return imageInfo, err
	}

	imageInfo.Width = img.Bounds().Max.X
	imageInfo.Height = img.Bounds().Max.Y
	imageInfo.BlockWidth = int(math.Floor(float64(imageInfo.Width) / 15))
	imageInfo.BlockHeight = int(math.Floor(float64(imageInfo.Height) / 15))

	for a := 0; a < int(imageInfo.Width/imageInfo.BlockWidth); a++ {
		for b := 0; b < int(imageInfo.Height/imageInfo.BlockHeight); b++ {
			x := a * imageInfo.BlockWidth
			y := b * imageInfo.BlockHeight
			xx := math.Min(float64(x+imageInfo.BlockWidth), float64(imageInfo.Width))
			yy := math.Min(float64(y+imageInfo.BlockHeight), float64(imageInfo.Height))

			a, r, g, b := img.At(x+int(math.Round(float64(imageInfo.BlockWidth)/2)), y+int(math.Round(float64(imageInfo.BlockHeight)/2))).RGBA()
			colorInfo := processor.GetHexColor(color.RGBA{uint8(a >> 8), uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)})

			pixelRectangle := strings.Join([]string{colorInfo, strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(int(xx)), strconv.Itoa(int(yy))}, "-")
			imageInfo.PixelInfo = append(imageInfo.PixelInfo, pixelRectangle)
		}
	}

	return imageInfo, nil
}

// GetHexColor turns RGBA to HEX
func (processor *ImageProcessor) GetHexColor(color color.RGBA) string {
	return fmt.Sprintf("%02x%02x%02x", color.R, color.G, color.B)
}

// ParseHexColorFast turns RGBA to HEX
func (processor *ImageProcessor) ParseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errors.New("Pixel color information is incorrect")
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errors.New("Pixel color information is incorrect")
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errors.New("Pixel color information is incorrect")
	}
	return
}

// ExtractPixelInfo turns compressed pixelinfo string into struct slice
func (processor *ImageProcessor) ExtractPixelInfo(s string) ([]PixelRectangle, error) {
	var rectanglesReturn []PixelRectangle

	rectangles := strings.Split(s, "_")
	for i := 0; i < len(rectangles); i++ {
		pixelRectangle, err := processor.ExtractRectangleInfo(rectangles[i])
		if err != nil {
			return rectanglesReturn, err
		}
		rectanglesReturn = append(rectanglesReturn, pixelRectangle)
	}

	return rectanglesReturn, nil
}

// ExtractRectangleInfo from imploded string
func (processor *ImageProcessor) ExtractRectangleInfo(s string) (PixelRectangle, error) {
	pixelRectangle := PixelRectangle{}
	rectangleInfo := strings.Split(s, "-")
	x, err := strconv.Atoi(rectangleInfo[1])
	if err != nil {
		return pixelRectangle, errors.New("Pixel position information is incorrect")
	}

	y, err := strconv.Atoi(rectangleInfo[2])
	if err != nil {
		return pixelRectangle, errors.New("Pixel position information is incorrect")
	}

	a, err := strconv.Atoi(rectangleInfo[3])
	if err != nil {
		return pixelRectangle, errors.New("Pixel position information is incorrect")
	}

	b, err := strconv.Atoi(rectangleInfo[4])
	if err != nil {
		return pixelRectangle, errors.New("Pixel position information is incorrect")
	}

	rectangle := image.Rect(x, y, a, b)
	pixelRectangle = PixelRectangle{Color: rectangleInfo[0], Rectangle: rectangle}
	return pixelRectangle, nil
}
