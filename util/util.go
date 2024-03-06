package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	result := make([]rune, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD, VND}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email address
func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}
