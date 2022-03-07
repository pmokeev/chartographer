package utils

import "fmt"

type IdError struct {
	ID int
}

func (error *IdError) Error() string {
	return fmt.Sprintf("Image with %v id does not exist ", error.ID)
}
