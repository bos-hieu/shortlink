package utils

import (
	"fmt"
	"testing"
)

func TestGenUniqueValue(t *testing.T) {
	uniqueValuesMap := map[string]bool{}
	type testType struct {
		name string
	}

	var tests []testType
	for i := range 100 {
		tests = append(tests, testType{
			name: fmt.Sprintf("test_%d", i),
		})
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenUniqueValue(); uniqueValuesMap[got] {
				t.Errorf("GenUniqueValue() = %v, existed in the map", got)
			} else {
				uniqueValuesMap[got] = true
			}
		})
	}
}
