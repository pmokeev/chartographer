package conrtollers

import "pmokeev/chartographer/internal/services"

type ChartographerController interface {
}

type Controller struct {
	ChartographerController
}

func NewController(service *services.Service) *Controller {
	return &Controller{ChartographerController: NewChartController(service.ChartographerServicer)}
}
