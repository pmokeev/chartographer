package utils

import "fmt"

type ParamsError struct {
}

func (error *ParamsError) Error() string {
	return fmt.Sprintf("Invalid parameters for image")
}
