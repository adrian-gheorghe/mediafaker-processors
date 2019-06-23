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
	X int
	Y int
	W int
	H int
}

// ImageInfo reflects information required by mediafaker to recreate the file. This amounts to width height and pixel info
type ImageInfo struct {
	W  int `json:"Width"`
	H  int `json:"Height"`
	P  []string
	BW int
	BH int
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

	imageInfo.W = img.Bounds().Max.X
	imageInfo.H = img.Bounds().Max.Y
	imageInfo.BW = 1
	imageInfo.BH = 1
	minWidth := int(math.Max(math.Floor(float64(imageInfo.W)/10), 50))
	minHeight := int(math.Max(math.Floor(float64(imageInfo.H)/10), 50))

	for i := minWidth; i > 4; i-- {
		if int(imageInfo.W%i) == 0 {
			imageInfo.BW = i
			break
		}
	}
	for i := minHeight; i > 4; i-- {
		if int(imageInfo.H%i) == 0 {
			imageInfo.BH = i
			break
		}
	}

	for a := 0; a < int(imageInfo.W/imageInfo.BW); a++ {
		for b := 0; b < int(imageInfo.H/imageInfo.BH); b++ {
			x := a * imageInfo.BW
			y := b * imageInfo.BH
			a, r, g, b := img.At(x+int(math.Round(float64(imageInfo.BW)/2)), y+int(math.Round(float64(imageInfo.BH)/2))).RGBA()
			colorInfo := processor.GetHexColor(color.RGBA{uint8(a >> 8), uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)})

			pixelRectangle := strings.Join([]string{colorInfo, strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(x + imageInfo.BW), strconv.Itoa(y + imageInfo.BH)}, "-")
			imageInfo.P = append(imageInfo.P, pixelRectangle)
		}
	}

	return imageInfo, nil
}

// GetHexColor turns RGBA to HEX
func (processor *ImageProcessor) GetHexColor(color color.RGBA) string {
	return fmt.Sprintf("%02x%02x%02x", color.R, color.G, color.B)
}
