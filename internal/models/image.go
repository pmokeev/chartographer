package models

import "sync"

type Image struct {
	ID     int
	Width  int
	Height int
	Mux    sync.Mutex
}

func NewImage(id, width, height int) *Image {
	return &Image{
		ID:     id,
		Width:  width,
		Height: height}
}
