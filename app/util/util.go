package util

func RatioInPercent(x, y float64) float64 {
	if y == 0 {
		return 0
	}

	return x / y * 100
}