package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestHashPassword(t *testing.T) {
	assert.Equal(t,
		HashPassword(""),
		"faf5d6fe5060d5b9e3f28f33794464d7e924184ab4bc230be17765eff22be81f")

	assert.Equal(t,
		HashPassword("known password"),
		"8124cb98f6a9343e753dbd6e95245f6092fdd8095d55d5dfcce50c32d7c3b3f3")
}

func TestMin(t *testing.T) {
	assert.Equal(t, Min(), 0)
	assert.Equal(t, Min(1), 1)
	assert.Equal(t, Min(3, 2, 1), 1)
	assert.Equal(t, Min(-3, -2, 1), -3)
}

func TestGenerateUsername(t *testing.T) {
	assert.Equal(t, len(GenerateUsername()), 9)
}
