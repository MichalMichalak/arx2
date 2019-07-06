package conf

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWordNormalization(t *testing.T) {
	type metadata struct {
		name string
		in   string
		out  string
	}
	tests := []metadata{
		{"empty", "", ""},
		{"basic_3_letters", "aaa", "aaa"},
		{"basic_1_letter", "a", "a"},
		{"1_lower_1_upper", "aA", "a-a"},
		{"2_words", "aAaa", "a-aaa"},
		{"3_words", "aaaBbbCcc", "aaa-bbb-ccc"},
		{"3_uppercase", "ABC", "a-b-c"},
		{"last_uppercase", "aaaB", "aaa-b"},
		{"first_uppercase", "AaaaBb", "aaaa-bb"},
		{"already_normalized", "aa-bb-cc", "aa-bb-cc"},
		{"number_last", "aa0", "aa0"},
		{"number_last_2_words", "aaAa0", "aa-aa0"},
		{"number_before_word", "aa0Bb", "aa0-bb"},
		{"number_in_middle", "aa0bb", "aa0bb"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testWordNormalization(t, test.in, test.out)
		})
	}
}

func testWordNormalization(t *testing.T, input, want string) {
	got := normalizeWord(input)
	require.Equal(t, got, want)
}
