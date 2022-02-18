package conrtollers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pmokeev/chartographer/internal/services"
)

type ChartController struct {
	chartService services.ChartographerServicer
}

func NewChartController(chartService services.ChartographerServicer) *ChartController {
	return &ChartController{chartService: chartService}
}

func (chartController *ChartController) CreateBMP(context *gin.Context) {
	fmt.Println("Create BMP")
}

func (chartController *ChartController) UpdateBMP(context *gin.Context) {
	fmt.Println("Update BMP")
}

func (chartController *ChartController) GetPartBMP(context *gin.Context) {
	fmt.Println("Get part of BMP")
}

func (chartController *ChartController) DeleteBMP(context *gin.Context) {
	fmt.Println("Delete BMP")
}
