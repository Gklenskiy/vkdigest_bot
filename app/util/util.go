package util

// RatioInPercent return ration in percent
func RatioInPercent(x, y float64) float64 {
	if y == 0 {
		return 0
	}

	return x / y * 100
}
