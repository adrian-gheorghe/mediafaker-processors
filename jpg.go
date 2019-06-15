package processors

import (
	"image"
	"image/color"
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
	Width       int `json:"Width"`
	Height      int `json:"Height"`
	PixelInfo   []PixelRectangle
	BlockWidth  int
	BlockHeight int
}

// PixelRectangle reflects the size, position and color of a pixel rectangle
type PixelRectangle struct {
	Color     image.Uniform
	Rectangle image.Rectangle
}

// JpgProcessor structure
type JpgProcessor struct {
}

// Inspect Command
func (processor *JpgProcessor) Inspect(sourcePath string, sourceInfo os.FileInfo) (ImageInfo, error) {
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
	imageInfo.BlockWidth = 1
	imageInfo.BlockHeight = 1

	for i := int(math.Floor(float64(imageInfo.Width) / 15)); i > 4; i-- {
		if int(imageInfo.Width%i) == 0 {
			imageInfo.BlockWidth = i
			break
		}
	}
	for i := int(math.Floor(float64(imageInfo.Height) / 15)); i > 4; i-- {
		if int(imageInfo.Height%i) == 0 {
			imageInfo.BlockHeight = i
			break
		}
	}

	for a := 0; a < int(imageInfo.Width/imageInfo.BlockWidth); a++ {
		for b := 0; b < int(imageInfo.Height/imageInfo.BlockHeight); b++ {
			x := a * imageInfo.BlockWidth
			y := b * imageInfo.BlockHeight
			rectangle := image.Rect(x, y, x+imageInfo.BlockWidth, y+imageInfo.BlockHeight)
			a, r, g, b := img.At(x+int(math.Round(float64(imageInfo.BlockWidth)/2)), y+int(math.Round(float64(imageInfo.BlockHeight)/2))).RGBA()
			pixelRectangle := PixelRectangle{Color: image.Uniform{color.RGBA{uint8(a >> 8), uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}}, Rectangle: rectangle}
			imageInfo.PixelInfo = append(imageInfo.PixelInfo, pixelRectangle)
		}
	}

	return imageInfo, nil
}