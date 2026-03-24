package mso

import (
	"reflect"
	"testing"
)

func TestSplitCommaString(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty", "", []string{}},
		{"whitespace_only", "   \t  ", []string{}},
		{"single_token", "1/1", []string{"1/1"}},
		{"single_token_with_spaces", " 1/1 ", []string{"1/1"}},
		{"multiple_tokens", "1/1,1/2", []string{"1/1", "1/2"}},
		{"multiple_tokens_with_spaces", " 1/1 , 1/2 ", []string{"1/1", "1/2"}},
		{"drops_empty_tokens_leading_trailing_commas", ",1/1,1/2,", []string{"1/1", "1/2"}},
		{"drops_empty_tokens_repeated_commas", "1/1,,1/2", []string{"1/1", "1/2"}},
		{"all_empty_tokens", ",,,", []string{}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := splitCommaString(testCase.input)
			if !reflect.DeepEqual(got, testCase.expected) {
				t.Fatalf("splitCommaString(%q) = %#v, expected %#v", testCase.input, got, testCase.expected)
			}
		})
	}
}
