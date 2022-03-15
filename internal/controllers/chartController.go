package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/pmokeev/chartographer/internal/models"
	"github.com/pmokeev/chartographer/internal/services"
	"golang.org/x/image/bmp"
	"io"
	"net/http"
	"strconv"
)

type ChartController struct {
	chartService services.ChartographerServicer
}

func NewChartController(chartService services.ChartographerServicer) *ChartController {
	return &ChartController{chartService: chartService}
}

func (chartController *ChartController) CreateBMP(context *gin.Context) {
	width, widthOk := context.GetQuery("width")
	height, heightOk := context.GetQuery("height")
	if !widthOk || !heightOk {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	createdID, err := chartController.chartService.CreateBMP(widthInt, heightInt)
	if err != nil {
		switch err.(type) {
		case *models.ParamsError:
			context.AbortWithStatus(http.StatusBadRequest)
			return
		default:
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	context.JSON(http.StatusCreated, map[string]int{
		"id": createdID,
	})
}

func (chartController *ChartController) UpdateBMP(context *gin.Context) {
	imageID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	xPosition, xPositionOk := context.GetQuery("x")
	yPosition, yPositionOk := context.GetQuery("y")
	width, widthOk := context.GetQuery("width")
	height, heightOk := context.GetQuery("height")
	if !widthOk || !heightOk || !xPositionOk || !yPositionOk {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	xPositionInt, err := strconv.Atoi(xPosition)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	yPositionInt, err := strconv.Atoi(yPosition)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	receivedImage, _, err := context.Request.FormFile("upload")
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, receivedImage); err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = chartController.chartService.UpdateBMP(imageID, xPositionInt, yPositionInt, widthInt, heightInt, buffer.Bytes())

	if err != nil {
		switch err.(type) {
		case *models.ParamsError:
			context.AbortWithStatus(http.StatusBadRequest)
			return
		case *models.IdError:
			context.AbortWithStatus(http.StatusNotFound)
			return
		default:
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	context.AbortWithStatus(http.StatusOK)
}

func (chartController *ChartController) GetPartBMP(context *gin.Context) {
	imageID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	xPosition, xPositionOk := context.GetQuery("x")
	yPosition, yPositionOk := context.GetQuery("y")
	width, widthOk := context.GetQuery("width")
	height, heightOk := context.GetQuery("height")
	if !widthOk || !heightOk || !xPositionOk || !yPositionOk {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	xPositionInt, err := strconv.Atoi(xPosition)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	yPositionInt, err := strconv.Atoi(yPosition)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	image, err := chartController.chartService.GetPartBMP(imageID, xPositionInt, yPositionInt, widthInt, heightInt)
	if err != nil {
		switch err.(type) {
		case *models.ParamsError:
			context.AbortWithStatus(http.StatusBadRequest)
			return
		case *models.IdError:
			context.AbortWithStatus(http.StatusNotFound)
			return
		default:
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	context.Header("Content-Type", "image/bmp")
	context.Stream(func(w io.Writer) bool {
		context.Status(200)
		bmp.Encode(w, image)
		return false
	})
}

func (chartController *ChartController) DeleteBMP(context *gin.Context) {
	imageID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = chartController.chartService.DeleteBMP(imageID); err != nil {
		switch err.(type) {
		case *models.ParamsError:
			context.AbortWithStatus(http.StatusBadRequest)
			return
		case *models.IdError:
			context.AbortWithStatus(http.StatusNotFound)
			return
		default:
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	context.AbortWithStatus(http.StatusOK)
}
