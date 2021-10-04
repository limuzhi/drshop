package util

import (
	"crypto/rand"
	"math"
	"math/big"
)

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(defaultLetters))))
		b[i] = defaultLetters[n.Int64()]
	}
	return string(b)
}


func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n // TODO +0.5 是为了四舍五入，如果不希望这样去掉这个
}
