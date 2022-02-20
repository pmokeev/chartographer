package models

type Image struct {
	ID            int
	Width         int
	Height        int
	OwnershipFile chan bool
}

func NewImage(id, width, height int) *Image {
	return &Image{
		ID:            id,
		Width:         width,
		Height:        height,
		OwnershipFile: make(chan bool)}
}
