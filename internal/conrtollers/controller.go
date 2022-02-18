package conrtollers

import (
	"github.com/gin-gonic/gin"
	"pmokeev/chartographer/internal/services"
)

type ChartographerController interface {
	CreateBMP(context *gin.Context)
	UpdateBMP(context *gin.Context)
	GetPartBMP(context *gin.Context)
	DeleteBMP(context *gin.Context)
}

type Controller struct {
	ChartographerController
}

func NewController(service *services.Service) *Controller {
	return &Controller{ChartographerController: NewChartController(service.ChartographerServicer)}
}
