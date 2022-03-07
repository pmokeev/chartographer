package services

import (
	"github.com/pmokeev/chartographer/internal/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/bmp"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func isEqualImages(actualImage, expectedImage image.Image) bool {
	actualBounds := actualImage.Bounds()
	expectedBounds := expectedImage.Bounds()
	if actualBounds.Dx() != expectedBounds.Dx() || actualBounds.Dy() != expectedBounds.Dy() {
		return false
	}

	for x := 0; x < expectedBounds.Dx(); x++ {
		for y := 0; y < expectedBounds.Dy(); y++ {
			if actualImage.At(x, y) != expectedImage.At(x, y) {
				return false
			}
		}
	}

	return true
}

func TestChartService_CreateBMP(t *testing.T) {
	tests := []struct {
		testName   string
		width      int
		height     int
		expectedID int
	}{
		{
			testName:   "0_id_640x426",
			width:      640,
			height:     426,
			expectedID: 0,
		},
		{
			testName:   "1_id_1x1",
			width:      1,
			height:     1,
			expectedID: 1,
		},
		{
			testName:   "2_id_10x10",
			width:      10,
			height:     10,
			expectedID: 2,
		},
		{
			testName:   "3_id_12x12",
			width:      12,
			height:     12,
			expectedID: 3,
		},
		{
			testName:   "Invalid width",
			width:      -1,
			height:     1,
			expectedID: -1,
		},
		{
			testName:   "Invalid height",
			width:      1,
			height:     -1,
			expectedID: -1,
		},
		{
			testName:   "Invalid width and height",
			width:      1,
			height:     -1,
			expectedID: -1,
		},
		{
			testName:   "Too big height",
			width:      1,
			height:     50001,
			expectedID: -1,
		},
		{
			testName:   "Too big width",
			width:      20001,
			height:     1,
			expectedID: -1,
		},
	}

	pathToStorageFolder := "../../testData/createBMP/"
	currentService := NewService(pathToStorageFolder)

	for ind, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualId, err := currentService.CreateBMP(test.width, test.height)
			if err != nil {
				assert.Equal(t, -1, actualId)
				assert.True(t, test.width <= 0 || test.width > 20000 || test.height <= 0 || test.height > 50000)
				return
			}
			assert.Equal(t, test.expectedID, actualId)

			expectedFile, err := os.OpenFile(filepath.Join(pathToStorageFolder, "correct"+strconv.Itoa(ind)+".bmp"), os.O_RDONLY, 0777)
			assert.NoError(t, err)
			expectedImage, err := bmp.Decode(expectedFile)
			assert.NoError(t, err)
			err = expectedFile.Close()
			assert.NoError(t, err)

			actualFile, err := os.OpenFile(filepath.Join(pathToStorageFolder, strconv.Itoa(ind)+".bmp"), os.O_RDONLY, 0777)
			assert.NoError(t, err)
			actualImage, err := bmp.Decode(actualFile)
			assert.NoError(t, err)
			err = actualFile.Close()
			assert.NoError(t, err)

			assert.True(t, isEqualImages(actualImage, expectedImage))

			err = os.Remove(pathToStorageFolder + "/" + strconv.Itoa(ind) + ".bmp")
			assert.NoError(t, err)
		})
	}
}

func TestChartService_UpdateBMP(t *testing.T) {
	tests := []struct {
		testName  string
		id        int
		xPosition int
		yPosition int
		width     int
		height    int
	}{
		{
			testName:  "Zero x and y",
			id:        0,
			xPosition: 0,
			yPosition: 0,
			width:     124,
			height:    124,
		},
		{
			testName:  "Positive x and y",
			id:        0,
			xPosition: 62,
			yPosition: 62,
			width:     124,
			height:    124,
		},
		{
			testName:  "Positive x and zero y",
			id:        0,
			xPosition: 62,
			yPosition: 0,
			width:     124,
			height:    124,
		},
		{
			testName:  "Positive y and zero x",
			id:        0,
			xPosition: 0,
			yPosition: 62,
			width:     124,
			height:    124,
		},
		{
			testName:  "Negative x and y",
			id:        0,
			xPosition: -62,
			yPosition: -62,
			width:     124,
			height:    124,
		},
		{
			testName:  "Negative x and zero y",
			id:        0,
			xPosition: -62,
			yPosition: 0,
			width:     124,
			height:    124,
		},
		{
			testName:  "Negative y and zero x",
			id:        0,
			xPosition: 0,
			yPosition: -62,
			width:     124,
			height:    124,
		},
		{
			testName:  "Negative x and positive y",
			id:        0,
			xPosition: -62,
			yPosition: 62,
			width:     124,
			height:    124,
		},
		{
			testName:  "Negative y and positive x",
			id:        0,
			xPosition: 62,
			yPosition: -62,
			width:     124,
			height:    124,
		},
		{
			testName:  "Less negative x and y",
			id:        0,
			xPosition: -125,
			yPosition: -125,
			width:     124,
			height:    124,
		},
		{
			testName:  "More positive x and y",
			id:        0,
			xPosition: 125,
			yPosition: 125,
			width:     124,
			height:    124,
		},
		{
			testName:  "More positive x",
			id:        0,
			xPosition: 125,
			yPosition: 62,
			width:     124,
			height:    124,
		},
		{
			testName:  "Wrong ID",
			id:        -1,
			xPosition: 62,
			yPosition: 62,
			width:     124,
			height:    124,
		},
		{
			testName:  "Wrong width",
			id:        0,
			xPosition: 62,
			yPosition: 62,
			width:     -1,
			height:    124,
		},
		{
			testName:  "Wrong height",
			id:        0,
			xPosition: 62,
			yPosition: 62,
			width:     124,
			height:    -1,
		},
	}

	pathToStorageFolder := "../../testData/updateBMP/"
	data, err := ioutil.ReadFile("../../testData/common/testImage.bmp")
	assert.NoError(t, err)

	for ind, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			currentService := NewService(pathToStorageFolder)

			_, err := currentService.CreateBMP(124, 124)
			assert.NoError(t, err)
			defer os.Remove(pathToStorageFolder + "/" + strconv.Itoa(test.id) + ".bmp")

			err = currentService.UpdateBMP(test.id, test.xPosition, test.yPosition, test.width, test.height, data)
			if err != nil {
				assert.True(t, test.width <= 0 || test.height <= 0 || test.id < 0 || utils.Abs(test.xPosition) >= test.width || utils.Abs(test.yPosition) >= test.height)
				return
			}

			expectedFile, err := os.OpenFile(filepath.Join(pathToStorageFolder, "correct"+strconv.Itoa(ind)+".bmp"), os.O_RDONLY, 0777)
			assert.NoError(t, err)
			expectedImage, err := bmp.Decode(expectedFile)
			assert.NoError(t, err)
			err = expectedFile.Close()
			assert.NoError(t, err)

			actualFile, err := os.OpenFile(filepath.Join(pathToStorageFolder, strconv.Itoa(test.id)+".bmp"), os.O_RDONLY, 0777)
			assert.NoError(t, err)
			actualImage, err := bmp.Decode(actualFile)
			assert.NoError(t, err)
			err = actualFile.Close()
			assert.NoError(t, err)

			assert.True(t, isEqualImages(actualImage, expectedImage))
		})
	}
}

func TestChartService_UpdateBMP_Overlapping(t *testing.T) {
	pathToStorageFolder := "../../testData/updateBMP/"
	data, err := ioutil.ReadFile("../../testData/common/testImage.bmp")
	assert.NoError(t, err)

	currentService := NewService(pathToStorageFolder)
	_, err = currentService.CreateBMP(124, 124)
	assert.NoError(t, err)
	defer os.Remove(pathToStorageFolder + "/" + strconv.Itoa(0) + ".bmp")
	err = currentService.UpdateBMP(0, 0, 0, 124, 124, data)
	assert.NoError(t, err)
	err = currentService.UpdateBMP(0, 62, 62, 124, 124, data)
	assert.NoError(t, err)

	expectedFile, err := os.OpenFile(filepath.Join(pathToStorageFolder, "correct9.bmp"), os.O_RDONLY, 0777)
	assert.NoError(t, err)
	expectedImage, err := bmp.Decode(expectedFile)
	assert.NoError(t, err)
	err = expectedFile.Close()
	assert.NoError(t, err)

	actualFile, err := os.OpenFile(filepath.Join(pathToStorageFolder, strconv.Itoa(0)+".bmp"), os.O_RDONLY, 0777)
	assert.NoError(t, err)
	actualImage, err := bmp.Decode(actualFile)
	assert.NoError(t, err)
	err = actualFile.Close()
	assert.NoError(t, err)

	assert.True(t, isEqualImages(actualImage, expectedImage))
}

func TestChartService_GetPartBMP(t *testing.T) {
	tests := []struct {
		testName  string
		width     int
		height    int
		id        int
		xPosition int
		yPosition int
	}{
		{
			testName:  "Zero x and y",
			width:     124,
			height:    124,
			id:        0,
			xPosition: 0,
			yPosition: 0,
		},
		{
			testName:  "Positive x and y",
			width:     124,
			height:    124,
			id:        0,
			xPosition: 62,
			yPosition: 62,
		},
		{
			testName:  "Positive x and zero y",
			width:     124,
			height:    124,
			id:        0,
			xPosition: 62,
			yPosition: 0,
		},
		{
			testName:  "Positive y and zero x",
			width:     124,
			height:    124,
			id:        0,
			xPosition: 0,
			yPosition: 62,
		},
		{
			testName:  "Negative x and y",
			width:     124,
			height:    124,
			id:        0,
			xPosition: -62,
			yPosition: -62,
		},
		{
			testName:  "Negative x and zero y",
			width:     124,
			height:    124,
			id:        0,
			xPosition: -62,
			yPosition: 0,
		},
		{
			testName:  "Negative y and zero x",
			width:     124,
			height:    124,
			id:        0,
			xPosition: 0,
			yPosition: -62,
		},
		{
			testName:  "Negative x and positive y",
			width:     124,
			height:    124,
			id:        0,
			xPosition: -62,
			yPosition: 62,
		},
		{
			testName:  "Negative y and positive x",
			width:     124,
			height:    124,
			id:        0,
			xPosition: 62,
			yPosition: -62,
		},
		{
			testName:  "Less x and y then width and height",
			width:     124,
			height:    124,
			id:        0,
			xPosition: -125,
			yPosition: -125,
		},
		{
			testName:  "More x and y then width and height",
			width:     124,
			height:    124,
			id:        0,
			xPosition: 125,
			yPosition: 125,
		},
		{
			testName:  "Invalid id",
			width:     124,
			height:    124,
			id:        -1,
			xPosition: 125,
			yPosition: 125,
		},
		{
			testName:  "Width less than zero ",
			width:     -1,
			height:    124,
			id:        0,
			xPosition: 125,
			yPosition: 125,
		},
		{
			testName:  "Height less than zero ",
			width:     124,
			height:    -1,
			id:        0,
			xPosition: 125,
			yPosition: 125,
		},
		{
			testName:  "Width is greater than 5000",
			width:     5001,
			height:    124,
			id:        0,
			xPosition: 125,
			yPosition: 125,
		},
		{
			testName:  "Height is greater than 5000",
			width:     124,
			height:    5001,
			id:        0,
			xPosition: 125,
			yPosition: 125,
		},
	}

	pathToStorageFolder := "../../testData/getPartBMP/"
	currentService := NewService(pathToStorageFolder)
	_, err := currentService.CreateBMP(124, 124)
	assert.NoError(t, err)

	data, err := ioutil.ReadFile("../../testData/common/testImage.bmp")
	assert.NoError(t, err)
	err = currentService.UpdateBMP(0, 0, 0, 124, 124, data)
	assert.NoError(t, err)

	for ind, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualImage, err := currentService.GetPartBMP(test.id, test.xPosition, test.yPosition, test.width, test.height)
			if err != nil {
				assert.True(t, test.width <= 0 || test.height <= 0 || test.id < 0 || test.width > 5000 || test.height > 5000 || utils.Abs(test.xPosition) >= test.width || utils.Abs(test.yPosition) >= test.height)
				return
			}
			assert.NoError(t, err)

			expectedFile, err := os.OpenFile(filepath.Join(pathToStorageFolder, "correct"+strconv.Itoa(ind)+".bmp"), os.O_RDONLY, 0777)
			assert.NoError(t, err)
			expectedImage, err := bmp.Decode(expectedFile)
			assert.NoError(t, err)
			err = expectedFile.Close()
			assert.NoError(t, err)

			assert.True(t, isEqualImages(actualImage, expectedImage))
		})
	}

	err = currentService.DeleteBMP(0)
	assert.NoError(t, err)
}

func TestChartService_DeleteBMP(t *testing.T) {
	tests := []struct {
		testName string
		width    int
		height   int
		id       int
	}{
		{
			testName: "OK",
			width:    124,
			height:   124,
			id:       0,
		},
		{
			testName: "OK",
			width:    10,
			height:   10,
			id:       1,
		},
		{
			testName: "Invalid ID",
			width:    10,
			height:   10,
			id:       -1,
		},
	}

	pathToStorageFolder := "../../testData/common/"
	currentService := NewService(pathToStorageFolder)

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			_, err := currentService.CreateBMP(test.width, test.height)
			assert.NoError(t, err)

			err = currentService.DeleteBMP(test.id)
			if err != nil {
				assert.True(t, test.id < 0)
				_ = currentService.DeleteBMP(2)
				return
			}
			assert.NoError(t, err)
		})
	}
}
