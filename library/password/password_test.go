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
		{name: "TooShortWithUpperCase", password: "A", flags: N | LorU, minLen: 2, maxLen: 0, expected: false},
		{name: "TooShortWithLowerCase", password: "a", flags: N | LorU, minLen: 2, maxLen: 0, expected: false},
		{name: "TooShortWithNumber", password: "1", flags: N | LorU, minLen: 2, maxLen: 0, expected: false},
		{name: "NoNumber", password: "Aa", flags: N | LorU, minLen: 2, maxLen: 0, expected: false},
		{name: "ValidWithUpperAndNumber", password: "A1", flags: N | LorU, minLen: 2, maxLen: 0, expected: true},
		{name: "ValidWithLowerAndNumber", password: "a1", flags: N | LorU, minLen: 2, maxLen: 0, expected: true},
		{name: "ValidWithAllConditions", password: "Aa1", flags: N | LorU, minLen: 2, maxLen: 0, expected: true},
		{name: "ValidWithUpperAndLower", password: "AaBbCcDd", flags: L | U, minLen: 6, maxLen: 0, expected: true},
		{name: "ValidWithUpperAndNumber", password: "ABCD1234", flags: N | U, minLen: 6, maxLen: 0, expected: true},
		{name: "ValidWithUpperAndNumberOrLetter", password: "ABCD1234", flags: N | LorU | U, minLen: 6, maxLen: 0, expected: true},
		{name: "MissingLower", password: "ABCD1234", flags: N | LorU | L, minLen: 6, maxLen: 0, expected: false},
		{name: "ExceedsMaxLen", password: "abcd1234", flags: N | L, minLen: 6, maxLen: 7, expected: false},
		{name: "ValidWithinRange", password: "Aa123456", flags: N | LorU, minLen: 6, maxLen: 8, expected: true},
		{name: "ValidWithSpecialChars", password: "Aa@123456", flags: Flag(1<<6) | N | L | U | S, minLen: 6, maxLen: 16, expected: true},
		{name: "ContainsNonASCII", password: "Aa123456语言", flags: N | L | U, minLen: 6, maxLen: 0, expected: false},
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
