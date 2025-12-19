package utils

import "strconv"

func StrAtoi(str string) int64 {
	atoi, _ := strconv.ParseInt(str, 10, 64)
	return atoi
}
