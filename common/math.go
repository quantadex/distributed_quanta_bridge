package common

func Min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Min32(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}