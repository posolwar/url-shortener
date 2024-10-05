package random

import "testing"

func TestValidZeroValue(t *testing.T) {
	result := NewRandomString(0)
	if result != "" {
		t.Errorf("Expected empty string, but got: %s", result)
	}
}

func TestValidSize(t *testing.T) {
	var size uint

	size = 5

	result := NewRandomString(size)
	if uint(len(result)) != size {
		t.Errorf("Expected %d, but got: %d", size, len(result))
	}
}
