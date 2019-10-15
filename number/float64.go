package number

// HalfEvent ...
func HalfEvent(x float64) int64 {
	if x == 0.0 {
		return 0
	}

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
	if x == 0.0 {
		return 0.0
	}

	return x * 1000.0
}

// Shift6 ...
func Shift6(x float64) float64 {
	if x == 0.0 {
		return 0.0
	}
	return x * 1000000.0
}

// Back3 ...
func Back3(x int64) float64 {
	if x == 0.0 {
		return 0.0
	}
	return float64(x) / 1000.0
}

// Back6 ...
func Back6(x int64) float64 {
	if x == 0.0 {
		return 0.0
	}
	return float64(x) / 1000000.0
}

// ShiftAndHalfEven6 ...
func ShiftAndHalfEven6(x float64) int64 {
	if x == 0.0 {
		return 0
	}
	return HalfEvent(Shift6(x))
}

// ShiftAndHalfEven3 ...
func ShiftAndHalfEven3(x float64) int64 {
	if x == 0.0 {
		return 0
	}

	return HalfEvent(Shift3(x))
}
