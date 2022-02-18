package services

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"pmokeev/chartographer/internal/models"
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
	currentImage.Mux.Lock()
	chartService.imageMap[chartService.idCounter] = currentImage

	file, err := os.Create(chartService.pathToStorageFolder + "/" + strconv.Itoa(currentImage.ID) + ".bmp")
	if err != nil {
		delete(chartService.imageMap, chartService.idCounter)
		currentImage.Mux.Unlock()
		return 0, err
	}
	err = png.Encode(file, img)
	if err != nil {
		delete(chartService.imageMap, chartService.idCounter)
		currentImage.Mux.Unlock()
		return 0, err
	}

	currentImage.Mux.Unlock()
	chartService.idCounter++
	return currentImage.ID, nil
}
