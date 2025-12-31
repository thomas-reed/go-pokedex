package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanInput(t *testing.T) {
	tests := map[string]struct{
		input string
		expected []string
	}{
		"extra spaces": {
			input: "    hello   world   ",
			expected: []string{"hello", "world"},
		},
		"mixed case": {
			input: "HeLLo WORLD how Are yOU?",
			expected: []string{"hello", "world", "how", "are", "you?"},
		},
		"just spaces/empty": {
			input: "   ",
			expected: []string{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T){
			actual := cleanInput(test.input)
			diff := cmp.Diff(test.expected, actual)
			if diff != "" {
				t.Fatalf("%s", diff)
			}
		})
	}
}