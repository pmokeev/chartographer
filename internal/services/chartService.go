package services

import (
	"bytes"
	"fmt"
	"github.com/pmokeev/chartographer/internal/models"
	"github.com/pmokeev/chartographer/internal/utils"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

type ChartService struct {
	imageMap            map[int]*models.Image
	pathToStorageFolder string
	idCounter           int
}

func NewChartService(pathToStorageFolder string) *ChartService {
	return &ChartService{
		pathToStorageFolder: pathToStorageFolder,
		idCounter:           0,
		imageMap:            make(map[int]*models.Image, 0)}
}

func (chartService *ChartService) CreateBMP(width, height int) (int, error) {
	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.Black)
		}
	}

	file, err := os.Create(chartService.pathToStorageFolder + "/" + strconv.Itoa(chartService.idCounter) + ".bmp")
	if err != nil {
		return 0, err
	}
	err = png.Encode(file, img)
	if err != nil {
		return 0, err
	}
	if err = file.Close(); err != nil {
		return 0, err
	}

	currentImage := models.NewImage(chartService.idCounter, width, height)
	chartService.imageMap[chartService.idCounter] = currentImage
	chartService.idCounter++

	return currentImage.ID, nil
}

func (chartService *ChartService) UpdateBMP(id, xPosition, yPosition, width, height int, receivedImageFile multipart.File) error {
	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &utils.RemoveError{ID: id}
	}

	originalImageFile, _ := os.Open(chartService.pathToStorageFolder + "/" + strconv.Itoa(id) + ".bmp")
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, receivedImageFile); err != nil {
		return err
	}
	originalImage, _, _ := image.Decode(originalImageFile)
	receivedImage, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		fmt.Println("HERE")
		return err // create special struct for this error
	}
	changeableOriginalImage, _ := originalImage.(*image.RGBA)

	for i := xPosition; i <= width; i++ {
		for j := yPosition; j <= height; j++ {
			if i <= currentImage.Width && j <= currentImage.Height {
				changeableOriginalImage.Set(i, j, receivedImage.At(i, j))
			}
		}
	}

	err = png.Encode(originalImageFile, originalImage)
	if err != nil {
		return err
	}

	originalImageFile.Close()
	receivedImageFile.Close()

	return nil
}

func (chartService *ChartService) GetPartBMP(id, xPosition, yPosition, width, height int) error {
	fmt.Println("Service get part of BMP")
	return nil
}

func (chartService *ChartService) DeleteBMP(id int) error {
	_, ok := chartService.imageMap[id]
	if !ok {
		return &utils.RemoveError{ID: id}
	}
	if err := os.Remove(chartService.pathToStorageFolder + "/" + strconv.Itoa(id) + ".bmp"); err != nil {
		return err
	}
	delete(chartService.imageMap, id)

	return nil
}
