package utils

import "fmt"

func FormatCurrency(n int) string {
	f := float32(n) / float32(100)
	return fmt.Sprintf("$ %.2f", f)
}
