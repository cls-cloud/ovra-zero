package utils

func BoolToIntFlag(b bool) int {
	if b {
		return 0
	}
	return 1
}
