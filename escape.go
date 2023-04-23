package xql

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	dollar     = '$'
	underscore = '_'
)

func EscapeName(name string, quote rune) string {
	var needQuote bool
	var hasQuote bool
	var digits int

	chars := []rune(name)

	for _, c := range chars {
		if unicode.IsDigit(c) {
			digits++
		}

		// ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar, underscore)
		if unicode.IsDigit(c) || unicode.IsLetter(c) || c == dollar || c == underscore {
			continue
		}

		needQuote = true
		hasQuote = c == quote

		break
	}

	// Identifiers may begin with a digit but unless quoted may not consist solely of digits.
	if len(chars) == digits {
		needQuote = true
	}

	if hasQuote {
		// Identifier quote characters can be included within an identifier if you quote the identifier.
		// If the character to be included within the identifier is the same as that used to quote the identifier itself,
		// then you need to double the character.
		name = strings.ReplaceAll(name, string(quote), string([]rune{quote, quote}))
	}

	if needQuote {
		return fmt.Sprintf("%c%s%c", quote, name, quote)
	}

	return name
}
