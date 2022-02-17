package conrtollers

import "pmokeev/chartographer/internal/services"

type ChartController struct {
	chartService services.ChartographerServicer
}

func NewChartController(chartService services.ChartographerServicer) *ChartController {
	return &ChartController{chartService: chartService}
}
