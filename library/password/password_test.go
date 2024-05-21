package password

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		flags    Flag
		minLen   int
		maxLen   int
		expected bool
	}{
		{"TooShortWithUpperCase", "A", N | LorU, 2, 0, false},
		{"TooShortWithLowerCase", "a", N | LorU, 2, 0, false},
		{"TooShortWithNumber", "1", N | LorU, 2, 0, false},
		{"NoNumber", "Aa", N | LorU, 2, 0, false},
		{"ValidWithUpperAndNumber", "A1", N | LorU, 2, 0, true},
		{"ValidWithLowerAndNumber", "a1", N | LorU, 2, 0, true},
		{"ValidWithAllConditions", "Aa1", N | LorU, 2, 0, true},
		{"ValidWithUpperAndLower", "AaBbCcDd", L | U, 6, 0, true},
		{"ValidWithUpperAndNumber", "ABCD1234", N | U, 6, 0, true},
		{"ValidWithUpperAndNumberOrLetter", "ABCD1234", N | LorU | U, 6, 0, true},
		{"MissingLower", "ABCD1234", N | LorU | L, 6, 0, false},
		{"ExceedsMaxLen", "abcd1234", N | L, 6, 7, false},
		{"ValidWithinRange", "Aa123456", N | LorU, 6, 8, true},
		{"ValidWithSpecialChars", "Aa@123456", Flag(1<<6) | N | L | U | S, 6, 16, true},
		{"ContainsNonASCII", "Aa123456语言", N | L | U, 6, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result bool
			if tc.maxLen == 0 {
				result = CheckPassword(tc.password, tc.flags, tc.minLen)
			} else {
				result = CheckPassword(tc.password, tc.flags, tc.minLen, tc.maxLen)
			}
			if result != tc.expected {
				t.Errorf("CheckPassword(%q, %d, %d, %d) = %v; want %v",
					tc.password, tc.flags, tc.minLen, tc.maxLen, result, tc.expected)
			}
		})
	}
}
