package services

type ChartographerServicer interface {
}

type Service struct {
	ChartographerServicer
}

func NewService() *Service {
	return &Service{ChartographerServicer: NewChartService()}
}
