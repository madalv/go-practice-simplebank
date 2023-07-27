package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghikjlmnopqrtsuvwxyz"

// called automatically when package is inited
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i:= 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	curr := []string{"EUR", "USD", "CAD"}
	n := len(curr)
	return curr[rand.Intn(n)]
}