// Package utils contains miscellaneous helper functions and classes
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
)

// Min takes any number of int arguments and returns the lowest. If no arguments are given it
// returns 0.
func Min(vars ...int) int {
	if len(vars) == 0 {
		return 0
	}

	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

// HashPassword returns the salted, hex encoded sha256 sum of the given string.
func HashPassword(in string) string {
	const salt string = "giIWQEmN6za#"
	hash := sha256.New()
	hash.Write([]byte(in + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateUsername returns a random 9 character long username.
func GenerateUsername() string {
	return strconv.Itoa(rand.Intn(999999999-100000000) + 100000000)
}
