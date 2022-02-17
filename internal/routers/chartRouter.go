package routers

import (
	"pmokeev/chartographer/internal/conrtollers"
	"pmokeev/chartographer/internal/services"
)

type chartRouter struct {
	controller *conrtollers.Controller
}

func NewChartRouter(service *services.Service) *chartRouter {
	return &chartRouter{controller: conrtollers.NewController(service)}
}
