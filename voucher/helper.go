package voucher

import (
	"log/slog"
	"strconv"
)

// convertStringToIntArray helper function to convert a string to an array of integers. This assumes the
// string has already been validated as a proper voucher.
func convertStringToIntArray(s string) ([]int, []int) {
	buffer := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		c, err := strconv.Atoi(string(s[i]))
		if err != nil {
			slog.Error("can not convert char to int", "char", string(s[i]), "error", err, "string", s)
			return []int{}, []int{}
		}
		buffer[i] = c
	}
	return buffer[:5], buffer[5:]
}
