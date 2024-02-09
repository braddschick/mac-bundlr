package types

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/nfnt/resize"

	"macbuilder/utils"
)

type Icon struct {
	FilePath string
	Width    int
	Height   int
	Img      image.Image
}

var macIconSizes = []uint{16, 32, 64, 128, 256, 512}

func NewIcon(filepath string, width, height int) (*Icon, error) {
	// currently only working with PNG files
	if strings.HasSuffix(filepath, ".png") {
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
	return nil, errors.New("PNG is the only image format currently handled")
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
	// Create a new folder to save the resampled images
	err := os.MkdirAll(outputFolderPath, os.ModePerm)
	if err != nil {
		return err
	}
	for _, size := range macIconSizes {
		err = i.Resample(path.Join(outputFolderPath, "icon_"+fmt.Sprint(size)+"x"+fmt.Sprint(size)+".png"), size, size)
		if err != nil {
			return err
		}
		err = i.Resample(path.Join(outputFolderPath, "icon_"+fmt.Sprint(size)+"x"+fmt.Sprint(size)+"@2x.png"), size*2, size*2)
		if err != nil {
			return err
		}
	}
	return err
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
