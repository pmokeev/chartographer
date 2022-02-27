package services

import (
	"bytes"
	"github.com/pmokeev/chartographer/internal/models"
	"github.com/pmokeev/chartographer/internal/utils"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"io"
	"mime/multipart"
	"os"
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
	currentImage := models.NewImage(chartService.idCounter, width, height, chartService.pathToStorageFolder+"/storage/"+strconv.Itoa(chartService.idCounter)+".bmp", true)
	currentImage.Lock()
	defer currentImage.Unlock()
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
		return 0, err
	}
	err = bmp.Encode(file, img)
	if err != nil {
		return 0, err
	}
	if err = file.Close(); err != nil {
		return 0, err
	}

	return currentImage.ID, nil
}

func (chartService *ChartService) UpdateBMP(id, xPosition, yPosition, width, height int, receivedImage multipart.File) error {
	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &utils.RemoveError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()
	if !currentImage.IsExist {
		return &utils.RemoveError{ID: id}
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

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, receivedImage); err != nil {
		return err
	}
	receivedImageDecoded, err := bmp.Decode(bytes.NewReader(buffer.Bytes()))
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

func (chartService *ChartService) GetPartBMP(id, xPosition, yPosition, width, height int) (string, error) {
	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return "", &utils.RemoveError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()
	if !currentImage.IsExist {
		return "", &utils.RemoveError{ID: id}
	}

	originalImageFile, err := os.OpenFile(currentImage.Filepath, os.O_RDONLY, 0777)
	if err != nil {
		return "", err
	}
	originalImage, err := bmp.Decode(originalImageFile)
	if err != nil {
		return "", err
	}
	if err := originalImageFile.Close(); err != nil {
		return "", err
	}

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	image := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x+xPosition >= 0 && x+xPosition < currentImage.Width && y+yPosition >= 0 && y+yPosition < currentImage.Height {
				image.Set(x, y, originalImage.At(x+xPosition, y+yPosition))
			} else {
				image.Set(x, y, color.Black)
			}
		}
	}

	pathToSavedFile := chartService.pathToStorageFolder + "/download/" + strconv.Itoa(chartService.counterGet) + ".bmp"
	chartService.counterGet++
	file, err := os.Create(pathToSavedFile)
	if err != nil {
		return "", err
	}
	err = bmp.Encode(file, image)
	if err != nil {
		return "", err
	}
	if err = file.Close(); err != nil {
		return "", err
	}

	return pathToSavedFile, nil
}

func (chartService *ChartService) DeleteBMP(id int) error {
	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &utils.RemoveError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()
	if !currentImage.IsExist {
		return &utils.RemoveError{ID: id}
	}
	if err := os.Remove(currentImage.Filepath); err != nil {
		return err
	}
	delete(chartService.imageMap, id)
	currentImage.IsExist = false

	return nil
}
