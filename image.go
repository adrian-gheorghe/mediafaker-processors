package processors

import (
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
