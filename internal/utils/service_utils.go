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

func Abs(number int) int {
	if number < 0 {
		return -number
	}
	return number
}

func Min(firstNumber, secondNumber int) int {
	if firstNumber < secondNumber {
		return firstNumber
	}
	return secondNumber
}
