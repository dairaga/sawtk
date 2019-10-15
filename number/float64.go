package number

// HalfEvent ...
func HalfEvent(x float64) int64 {
	var sign int64 = 1
	if x < 0 {
		sign = -1
		x = -x
	}

	x1 := int64(x)
	x2 := int64(x * 10.0)
	r := x2 - x1*10

	if r >= 0 && r <= 4 {
		return x1 * sign
	} else if r >= 6 && r <= 9 {
		return (x1 + 1) * sign
	} else {
		if (x1 & 0x01) == 0 {
			return x1 * sign
		}

		return (x1 + 1) * sign
	}
}

// Shift3 ...
func Shift3(x float64) float64 {
	return x * 1000.0
}

// Shift6 ...
func Shift6(x float64) float64 {
	return x * 1000000.0
}
