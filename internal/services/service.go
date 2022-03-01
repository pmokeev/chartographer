package services

import "image"

//go:generate mockgen -source=service.go -destination=./mocks/mock.go

type ChartographerServicer interface {
	CreateBMP(width, height int) (int, error)
	UpdateBMP(id, xPosition, yPosition, width, height int, receivedImage []byte) error
	GetPartBMP(id, xPosition, yPosition, width, height int) (image.Image, error)
	DeleteBMP(id int) error
}

type Service struct {
	ChartographerServicer
}

func NewService(pathToStorageFolder string) *Service {
	return &Service{ChartographerServicer: NewChartService(pathToStorageFolder)}
}
