package models

type Image struct {
	ID       int
	Width    int
	Height   int
	Filepath string
	IsExist  bool
}

func NewImage(id, width, height int, filepath string, isExist bool) *Image {
	return &Image{
		ID:       id,
		Width:    width,
		Height:   height,
		Filepath: filepath,
		IsExist:  isExist}
}
