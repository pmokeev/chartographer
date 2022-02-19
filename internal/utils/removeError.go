package utils

import "fmt"

type RemoveError struct {
	ID int
}

func (error *RemoveError) Error() string {
	return fmt.Sprintf("Image with %v id does not exist ", error.ID)
}
