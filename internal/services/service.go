package services

type ChartographerServicer interface {
	CreateBMP(width, height int) (int, error)
	UpdateBMP(id, xPosition, yPosition, width, height int) error
	GetPartBMP(id, xPosition, yPosition, width, height int) error
	DeleteBMP(id int) error
}

type Service struct {
	ChartographerServicer
}

func NewService(pathToStorageFolder string) *Service {
	return &Service{ChartographerServicer: NewChartService(pathToStorageFolder)}
}
