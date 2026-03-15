package analyzer

import (
	"unicode"
	"unicode/utf8"
)

func isLowercase(msg string) bool {
	if len(msg) == 0 {
		return true
	}
	r, _ := utf8.DecodeRuneInString(msg)
	return unicode.IsLower(r)
}

func containsEmoji(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.S, r) {
			return true
		}
		// Additional emoji ranges.
		if (r >= 0x1F600 && r <= 0x1F64F) || // emoticons
			(r >= 0x1F300 && r <= 0x1F5FF) || // symbols and pictographs
			(r >= 0x1F680 && r <= 0x1F6FF) { // transport and symbols
			return true
		}
	}
	return false
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

func containsSpecialChars(s string) bool {
	for _, r := range s {
		if r > 0x7F {
			continue
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' {
			continue
		}
		return true
	}
	return false
}
