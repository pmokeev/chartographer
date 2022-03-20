package services

import (
	"bytes"
	"github.com/pmokeev/chartographer/internal/models"
	"github.com/pmokeev/chartographer/internal/utils"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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
		return -1, &models.ParamsError{}
	}

	currentImage := models.NewImage(chartService.idCounter, width, height, filepath.Join(chartService.pathToStorageFolder, strconv.Itoa(chartService.idCounter)+".bmp"), true)
	currentImage.Lock()
	defer currentImage.Unlock()

	chartService.imageMap[chartService.idCounter] = currentImage
	chartService.idCounter++

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup
	goroutineCount := 10
	chunkSize := width/goroutineCount + 1
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)

		go func(image *image.RGBA, width, height, chunkSize, chunkNumber int) {
			defer wg.Done()
			for x := chunkSize * chunkNumber; x < utils.Min((chunkNumber+1)*chunkSize, width); x++ {
				for y := 0; y < height; y++ {
					image.Set(x, y, color.Black)
				}
			}
		}(img, width, height, chunkSize, i)
	}
	wg.Wait()

	file, err := os.Create(currentImage.Filepath)
	if err != nil {
		return 0, err
	}
	if err := utils.WriteInFile(file, img); err != nil {
		return 0, err
	}

	return currentImage.ID, nil
}

func (chartService *ChartService) UpdateBMP(id, xPosition, yPosition, width, height int, receivedImage []byte) error {
	if width <= 0 || height <= 0 {
		return &models.ParamsError{}
	}

	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &models.IdError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()

	if !currentImage.IsExist {
		return &models.IdError{ID: id}
	}

	if utils.Abs(xPosition) >= currentImage.Width || utils.Abs(yPosition) >= currentImage.Height {
		return &models.ParamsError{}
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

	changeableOriginalImage := image.NewRGBA(originalImage.Bounds())
	draw.Draw(changeableOriginalImage, originalImage.Bounds(), originalImage, image.Point{}, draw.Over)

	receivedImageDecoded, err := bmp.Decode(bytes.NewReader(receivedImage))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	goroutineCount := 10
	chunkSize := width/goroutineCount + 1
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)

		go func(image *image.RGBA, width, height, chunkSize, chunkNumber, xPosition, yPosition int) {
			defer wg.Done()
			for x := chunkSize * chunkNumber; x < utils.Min((chunkNumber+1)*chunkSize, width); x++ {
				for y := 0; y < height; y++ {
					if x+xPosition >= 0 && x+xPosition < currentImage.Width && y+yPosition >= 0 && y+yPosition < currentImage.Height {
						image.Set(x+xPosition, y+yPosition, receivedImageDecoded.At(x, y))
					}
				}
			}
		}(changeableOriginalImage, width, height, chunkSize, i, xPosition, yPosition)
	}
	wg.Wait()

	originalImageFile, err = os.OpenFile(currentImage.Filepath, os.O_WRONLY, 0777)
	if err := utils.WriteInFile(originalImageFile, changeableOriginalImage); err != nil {
		return err
	}

	return nil
}

func (chartService *ChartService) GetPartBMP(id, xPosition, yPosition, width, height int) (image.Image, error) {
	if width <= 0 || height <= 0 || width > 5000 || height > 5000 {
		return nil, &models.ParamsError{}
	}

	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return nil, &models.IdError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()

	if !currentImage.IsExist {
		return nil, &models.IdError{ID: id}
	}

	if utils.Abs(xPosition) >= currentImage.Width || utils.Abs(yPosition) >= currentImage.Height {
		return nil, &models.ParamsError{}
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

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup
	goroutineCount := 10
	chunkSize := width/goroutineCount + 1
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)

		go func(image *image.RGBA, width, height, chunkSize, chunkNumber, xPosition, yPosition int) {
			defer wg.Done()
			for x := chunkSize * chunkNumber; x < utils.Min((chunkNumber+1)*chunkSize, width); x++ {
				for y := 0; y < height; y++ {
					if x+xPosition >= 0 && x+xPosition < currentImage.Width && y+yPosition >= 0 && y+yPosition < currentImage.Height {
						image.Set(x, y, originalImage.At(x+xPosition, y+yPosition))
					} else {
						image.Set(x, y, color.Black)
					}
				}
			}
		}(img, width, height, chunkSize, i, xPosition, yPosition)
	}
	wg.Wait()

	return img, nil
}

func (chartService *ChartService) DeleteBMP(id int) error {
	currentImage, ok := chartService.imageMap[id]
	if !ok {
		return &models.IdError{ID: id}
	}
	currentImage.Lock()
	defer currentImage.Unlock()
	if !currentImage.IsExist {
		return &models.IdError{ID: id}
	}
	if err := os.Remove(currentImage.Filepath); err != nil {
		return err
	}
	delete(chartService.imageMap, id)
	currentImage.IsExist = false

	return nil
}
