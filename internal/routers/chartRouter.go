package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/pmokeev/chartographer/internal/conrtollers"
	"github.com/pmokeev/chartographer/internal/services"
)

type ChartRouter struct {
	controller *conrtollers.Controller
}

func NewChartRouter(service *services.Service) *ChartRouter {
	return &ChartRouter{controller: conrtollers.NewController(service)}
}

func (chartRouter *ChartRouter) InitChartRouter() *gin.Engine {
	router := gin.New()

	chart := router.Group("/chartas")
	{
		chart.POST("/", chartRouter.controller.CreateBMP)
		chart.POST("/:id/", chartRouter.controller.UpdateBMP)
		chart.GET("/:id/", chartRouter.controller.GetPartBMP)
		chart.DELETE("/:id/", chartRouter.controller.DeleteBMP)
	}

	return router
}
