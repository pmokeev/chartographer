package services

import "mime/multipart"

//go:generate mockgen -source=service.go -destination=../../tests/mocks/mock.go

type ChartographerServicer interface {
	CreateBMP(width, height int) (int, error)
	UpdateBMP(id, xPosition, yPosition, width, height int, receivedImage multipart.File) error
	GetPartBMP(id, xPosition, yPosition, width, height int) (string, error)
	DeleteBMP(id int) error
}

type Service struct {
	ChartographerServicer
}

func NewService(pathToStorageFolder string) *Service {
	return &Service{ChartographerServicer: NewChartService(pathToStorageFolder)}
}
