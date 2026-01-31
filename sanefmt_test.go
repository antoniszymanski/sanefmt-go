// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package sanefmt

import (
	"testing"
	"unicode/utf8"
)

func TestReplaceFirstRune(t *testing.T) {
	t.Run("unchangedRune", testReplaceFirstRune("aaaa", 'a'))
	t.Run("sameSize", testReplaceFirstRune("aaaa", 'b'))
	t.Run("smallerRune", testReplaceFirstRune("aaaa", 'ą'))
	t.Run("largerRune", testReplaceFirstRune("ąaaa", 'a'))
}

func FuzzReplaceFirstRune(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string, r rune) {
		testReplaceFirstRune(s, r)(t)
	})
}

func testReplaceFirstRune(s string, r rune) func(t *testing.T) {
	return func(t *testing.T) {
		expected := func() string {
			_, size := utf8.DecodeRuneInString(s)
			if size == 0 {
				return ""
			}
			return string(r) + s[size:]
		}()
		actual := string(replaceFirstRune([]byte(s), func(_ rune) rune { return r }))
		if actual != expected {
			t.Errorf("replaceFirstRune(%q, %q) = %q, should be %q", s, r, actual, expected)
		}
	}
}
