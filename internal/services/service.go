package services

type ChartographerServicer interface {
	CreateBMP(width, height int) (int, error)
}

type Service struct {
	ChartographerServicer
}

func NewService(pathToStorageFolder string) *Service {
	return &Service{ChartographerServicer: NewChartService(pathToStorageFolder)}
}
