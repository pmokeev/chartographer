package services

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

type ChartService struct {
	pathToStorageFolder string
	idCounter           int
}

func NewChartService(pathToStorageFolder string) *ChartService {
	return &ChartService{
		pathToStorageFolder: pathToStorageFolder,
		idCounter:           -1}
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

	chartService.idCounter++
	file, err := os.Create(chartService.pathToStorageFolder + "/" + strconv.Itoa(chartService.idCounter) + ".png")
	if err != nil {
		chartService.idCounter--
		return 0, err
	}
	err = png.Encode(file, img)
	if err != nil {
		chartService.idCounter--
		return 0, err
	}

	return chartService.idCounter, nil
}
