package services

import (
	"bytes"
	"github.com/pmokeev/chartographer/internal/models"
	"github.com/pmokeev/chartographer/internal/utils"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
)

type ChartService struct {
	imageMap            map[int]*models.Image
	pathToStorageFolder string
	idCounter           int
	counterGet          int
}

func NewChartService(pathToStorageFolder string) *ChartService {
	return &ChartService{
		pathToStorageFolder: pathToStorageFolder,
		idCounter:           0,
		imageMap:            make(map[int]*models.Image, 0),
		counterGet:          0}
}

func (chartService *ChartService) CreateBMP(width, height int) (int, error) {
	if width <= 0 || width > 20000 || height <= 0 || height > 50000 {
		return -1, &utils.ParamsError{}
	}

	currentImage := models.NewImage(chartService.idCounter, width, height, filepath.Join(chartService.pathToStorageFolder, strconv.Itoa(chartService.idCounter)+".bmp"), true)
	currentImage.Lock()
	defer currentImage.Unlock()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.Black)
		}
	}

	file, err := os.Create(currentImage.Filepath)
	if err != nil {
		return 0, err
	}
	err = bmp.Encode(file, img)
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

func (chartService *ChartService) UpdateBMP(id, xPosition, yPosition, width, height int, receivedImage []byte) error {
	if width <= 0 || height <= 0 {
		return &utils.ParamsError{}
	}

	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &utils.IdError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()
	if !currentImage.IsExist {
		return &utils.IdError{ID: id}
	}

	if utils.Abs(xPosition) >= currentImage.Width || utils.Abs(yPosition) >= currentImage.Height {
		return &utils.ParamsError{}
	}

	originalImageFile, err := os.OpenFile(currentImage.Filepath, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	originalImage, err := bmp.Decode(originalImageFile)
	if err != nil {
		return err
	}
	if err = originalImageFile.Close(); err != nil {
		return err
	}
	changeableOriginalImage, _ := originalImage.(*image.RGBA)

	receivedImageDecoded, err := bmp.Decode(bytes.NewReader(receivedImage))
	if err != nil {
		return err
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x+xPosition >= 0 && x+xPosition < currentImage.Width && y+yPosition >= 0 && y+yPosition < currentImage.Height {
				changeableOriginalImage.Set(x+xPosition, y+yPosition, receivedImageDecoded.At(x, y))
			}
		}
	}

	originalImageFile, err = os.OpenFile(currentImage.Filepath, os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	err = bmp.Encode(originalImageFile, changeableOriginalImage)
	if err != nil {
		return err
	}
	if err = originalImageFile.Close(); err != nil {
		return err
	}

	return nil
}

func (chartService *ChartService) GetPartBMP(id, xPosition, yPosition, width, height int) (image.Image, error) {
	if width <= 0 || height <= 0 || width > 5000 || height > 5000 {
		return nil, &utils.ParamsError{}
	}

	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return nil, &utils.IdError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()
	if !currentImage.IsExist {
		return nil, &utils.IdError{ID: id}
	}

	if utils.Abs(xPosition) >= currentImage.Width || utils.Abs(yPosition) >= currentImage.Height {
		return nil, &utils.ParamsError{}
	}

	originalImageFile, err := os.OpenFile(currentImage.Filepath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	originalImage, err := bmp.Decode(originalImageFile)
	if err != nil {
		return nil, err
	}
	if err := originalImageFile.Close(); err != nil {
		return nil, err
	}

	image := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x+xPosition >= 0 && x+xPosition < currentImage.Width && y+yPosition >= 0 && y+yPosition < currentImage.Height {
				image.Set(x, y, originalImage.At(x+xPosition, y+yPosition))
			} else {
				image.Set(x, y, color.Black)
			}
		}
	}

	return image, nil
}

func (chartService *ChartService) DeleteBMP(id int) error {
	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &utils.IdError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()
	if !currentImage.IsExist {
		return &utils.IdError{ID: id}
	}
	if err := os.Remove(currentImage.Filepath); err != nil {
		return err
	}
	delete(chartService.imageMap, id)
	currentImage.IsExist = false

	return nil
}
