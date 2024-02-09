package types

import (
	"image"
	"image/png"
	"os"
	"path"

	"github.com/jackmordaunt/icns"
	"github.com/nfnt/resize"

	"mac-bundlr/utils"
)

type Icon struct {
	FilePath string
	Width    int
	Height   int
	Img      image.Image
}

func NewIcon(filepath string, width, height int) (*Icon, error) {
	img, err := utils.CreateImg(filepath)
	if err != nil {
		return nil, err
	}
	return &Icon{
		FilePath: filepath,
		Width:    width,
		Height:   height,
		Img:      img,
	}, nil
}

func (i *Icon) Resample(outputFilePath string, width, height uint) error {
	err := utils.CopyFile(i.FilePath, outputFilePath)
	if err != nil {
		return err
	}
	img, err := utils.CreateImg(outputFilePath)
	if err != nil {
		return err
	}
	// Resample the image to the desired dimensions
	resampledImg := size(img, width, height)

	// Create a new file to save the resampled image
	resampledFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer resampledFile.Close()

	// Encode the resampled image as PNG and save it to the file
	err = png.Encode(resampledFile, resampledImg)
	if err != nil {
		return err
	}

	return nil
}

func (i *Icon) CreateMacIcons(outputFolderPath string) error {
	// Create a new file for the ICNS encoder
	out := path.Join(outputFolderPath, "icon.icns")
	file, err := os.Create(out)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new ICNS encoder
	encoder := icns.NewEncoder(file)

	// Define the icon sizes
	sizes := []uint{16, 32, 64, 128, 256, 512, 1024}

	// Add each size to the encoder
	for _, size := range sizes {
		resized := resize.Resize(size, size, i.Img, resize.Lanczos3)
		if err := encoder.Encode(resized); err != nil {
			return err
		}
	}
	return nil
}

func (i *Icon) Exists() bool {
	_, err := os.Stat(i.FilePath)
	return !os.IsNotExist(err)
}

func size(img image.Image, width, height uint) image.Image {
	// Resize the image to the desired dimensions using the Nearest Neighbor algorithm
	resampledImg := resize.Resize(width, height, img, resize.NearestNeighbor)
	return resampledImg
}
