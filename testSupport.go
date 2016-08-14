package main

import (
	"testing"
	"strconv"
)

func expect(t *testing.T, valueName string, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected %v to be %v but was %v", valueName, expected, actual)
	}
}

func expectInt(t *testing.T, valueName string, expected int, actual int) {
	expect(t, valueName, strconv.Itoa(expected), strconv.Itoa(actual))
}