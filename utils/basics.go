package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func CreateImg(src string) (image.Image, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	o := strings.ToLower(src)
	if strings.HasSuffix(o, ".png") {
		return png.Decode(file)
	}
	if strings.HasSuffix(o, ".jpg") || strings.HasSuffix(o, ".jpeg") {
		return jpeg.Decode(file)
	}
	return nil, errors.New("only PNG and JPEG files are currently supported")
}
