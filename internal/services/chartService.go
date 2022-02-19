package services

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"pmokeev/chartographer/internal/models"
	"pmokeev/chartographer/internal/utils"
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

	currentImage := models.NewImage(chartService.idCounter, width, height)

	file, err := os.Create(chartService.pathToStorageFolder + "/" + strconv.Itoa(currentImage.ID) + ".bmp")
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

	chartService.imageMap[chartService.idCounter] = currentImage
	chartService.idCounter++
	return currentImage.ID, nil
}

func (chartService *ChartService) UpdateBMP(id, xPosition, yPosition, width, height int) error {
	fmt.Println("Service update BMP")
	return nil
}

func (chartService *ChartService) GetPartBMP(id, xPosition, yPosition, width, height int) error {
	fmt.Println("Service get part of BMP")
	return nil
}

func (chartService *ChartService) DeleteBMP(id int) error {
	if _, ok := chartService.imageMap[id]; !ok {
		return &utils.RemoveError{ID: id}
	}
	delete(chartService.imageMap, id)
	if err := os.Remove(chartService.pathToStorageFolder + "/" + strconv.Itoa(id) + ".bmp"); err != nil {
		return err
	}

	return nil
}
