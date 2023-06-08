package nwc24

import (
	"fmt"
	"testing"
)

func TestLoadWiiNumber(t *testing.T) {
	wiiNumber := LoadWiiNumber(7615178254839298)
	fmt.Println(wiiNumber.GetHollywoodID())
}

func TestWiiNumber_CheckWiiNumber(t *testing.T) {
	// Erase any digit and it will throw
	wiiNumber := LoadWiiNumber(7615178254839298)
	if !wiiNumber.CheckWiiNumber() {
		t.Errorf("wii number is not valid")
	}
}
