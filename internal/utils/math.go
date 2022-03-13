package utils

func Abs(number int) int {
	if number < 0 {
		return -number
	}
	return number
}

func Min(firstNumber, secondNumber int) int {
	if firstNumber < secondNumber {
		return firstNumber
	}
	return secondNumber
}
