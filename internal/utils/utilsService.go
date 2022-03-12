package utils

import (
	"golang.org/x/image/bmp"
	"image"
	"os"
)

func WriteInFile(file *os.File, image image.Image) error {
	err := bmp.Encode(file, image)
	if err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}

	return nil
}
