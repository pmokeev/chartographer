package models

import "sync"

type Image struct {
	ID       int
	Width    int
	Height   int
	Filepath string
	IsExist  bool

	sync.Mutex
}

func NewImage(id, width, height int, filepath string, isExist bool) *Image {
	return &Image{
		ID:       id,
		Width:    width,
		Height:   height,
		Filepath: filepath,
		IsExist:  isExist}
}
