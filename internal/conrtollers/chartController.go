package conrtollers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pmokeev/chartographer/internal/services"
	"pmokeev/chartographer/internal/utils"
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
	if widthInt <= 0 || widthInt > 20000 || heightInt <= 0 || heightInt > 50000 {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	createdID, err := chartController.chartService.CreateBMP(widthInt, heightInt)
	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusCreated, map[string]int{
		"id": createdID,
	})
}

func (chartController *ChartController) UpdateBMP(context *gin.Context) {
	fmt.Println("Update BMP")
}

func (chartController *ChartController) GetPartBMP(context *gin.Context) {
	fmt.Println("Get part of BMP")
}

func (chartController *ChartController) DeleteBMP(context *gin.Context) {
	imageID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = chartController.chartService.DeleteBMP(imageID); err != nil {
		switch err.(type) {
		case *utils.RemoveError:
			context.AbortWithStatus(http.StatusBadRequest)
			return
		default:
			fmt.Println(err.Error())
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	context.AbortWithStatus(http.StatusOK)
}
