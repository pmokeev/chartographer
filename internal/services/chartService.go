package services

import (
	"bytes"
	"fmt"
	"github.com/pmokeev/chartographer/internal/models"
	"github.com/pmokeev/chartographer/internal/utils"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"sync"
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
	var mutex sync.Mutex
	currentImage := models.NewImage(chartService.idCounter, width, height, chartService.pathToStorageFolder+"/"+strconv.Itoa(chartService.idCounter)+".bmp", true)
	mutex.Lock()
	chartService.imageMap[chartService.idCounter] = currentImage
	chartService.idCounter++

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.Black)
		}
	}

	file, err := os.Create(currentImage.Filepath)
	if err != nil {
		mutex.Unlock()
		return 0, err
	}
	err = bmp.Encode(file, img)
	if err != nil {
		mutex.Unlock()
		return 0, err
	}
	if err = file.Close(); err != nil {
		mutex.Unlock()
		return 0, err
	}

	mutex.Unlock()
	return currentImage.ID, nil
}

func (chartService *ChartService) UpdateBMP(id, xPosition, yPosition, width, height int, receivedImage multipart.File) error {
	var mutex sync.Mutex
	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &utils.RemoveError{ID: id}
	}
	mutex.Lock()
	if !currentImage.IsExist {
		mutex.Unlock()
		return &utils.RemoveError{ID: id}
	}

	originalImageFile, _ := os.OpenFile(currentImage.Filepath, os.O_RDONLY, 0777)
	originalImage, _ := bmp.Decode(originalImageFile)
	originalImageFile.Close()
	changeableOriginalImage, _ := originalImage.(*image.RGBA)

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, receivedImage); err != nil {
		mutex.Unlock()
		return err
	}
	receivedImageDecoded, err := bmp.Decode(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		mutex.Unlock()
		return err
	}

	for x := xPosition; x <= width; x++ {
		for y := yPosition; y <= height; y++ {
			if x <= currentImage.Width && y <= currentImage.Height {
				changeableOriginalImage.Set(x, y, receivedImageDecoded.At(x-xPosition, y-yPosition))
			}
		}
	}

	originalImageFile, _ = os.OpenFile(currentImage.Filepath, os.O_WRONLY, 0777)
	err = bmp.Encode(originalImageFile, changeableOriginalImage)
	if err != nil {
		mutex.Unlock()
		return err
	}
	originalImageFile.Close()
	mutex.Unlock()

	return nil
}

func (chartService *ChartService) GetPartBMP(id, xPosition, yPosition, width, height int) error {
	fmt.Println("Service get part of BMP")
	return nil
}

func (chartService *ChartService) DeleteBMP(id int) error {
	var mutex sync.Mutex
	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &utils.RemoveError{ID: id}
	}
	mutex.Lock()
	if !currentImage.IsExist {
		mutex.Unlock()
		return &utils.RemoveError{ID: id}
	}
	if err := os.Remove(currentImage.Filepath); err != nil {
		mutex.Unlock()
		return err
	}
	delete(chartService.imageMap, id)
	currentImage.IsExist = false
	mutex.Unlock()

	return nil
}
