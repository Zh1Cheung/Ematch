package Ematch

import "unicode/utf8"

func isMatch(str, pattern string) bool {
	if pattern == "*" {
		return true
	}
	return wildcardMatch(str, pattern)
}

func IsPattern(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == '*' || str[i] == '?' {
			return true
		}
	}
	return false
}

func wildcardMatch(str, pattern string) bool {
	for len(pattern) > 0 {
		if pattern[0] > 0x7f {
			return boundaryProcessForCode(str, pattern)
		}
		switch pattern[0] {
		default:
			if len(str) == 0 {
				return false
			}
			if str[0] > 0x7f {
				return boundaryProcessForCode(str, pattern)
			}
			if str[0] != pattern[0] {
				return false
			}
		case '?':
			if len(str) == 0 {
				return false
			}
		case '*':
			return wildcardMatch(str, pattern[1:]) || (len(str) > 0 && wildcardMatch(str[1:], pattern))
		}
		str = str[1:]
		pattern = pattern[1:]
	}
	return len(str) == 0 && len(pattern) == 0
}

func boundaryProcessForCode(str, pattern string) bool {
	var st, pa rune
	var stria, patterns int

	if len(str) > 0 {
		if str[0] > 0x7f {
			st, stria = utf8.DecodeRuneInString(str)
		} else {
			st, stria = rune(str[0]), 1
		}
	} else {
		st, stria = utf8.RuneError, 0
	}
	if len(pattern) > 0 {
		if pattern[0] > 0x7f {
			pa, patterns = utf8.DecodeRuneInString(pattern)
		} else {
			pa, patterns = rune(pattern[0]), 1
		}
	} else {
		pa, patterns = utf8.RuneError, 0
	}

	for pa != utf8.RuneError {
		switch pa {
		default:
			if stria == utf8.RuneError {
				return false
			}
			if st != pa {
				return false
			}
		case '?':
			if stria == utf8.RuneError {
				return false
			}
		case '*':
			return boundaryProcessForCode(str, pattern[patterns:]) ||
				(stria > 0 && boundaryProcessForCode(str[stria:], pattern))
		}
		str = str[stria:]
		pattern = pattern[patterns:]

		if len(str) > 0 {
			if str[0] > 0x7f {
				st, stria = utf8.DecodeRuneInString(str)
			} else {
				st, stria = rune(str[0]), 1
			}
		} else {
			st, stria = utf8.RuneError, 0
		}
		if len(pattern) > 0 {
			if pattern[0] > 0x7f {
				pa, patterns = utf8.DecodeRuneInString(pattern)
			} else {
				pa, patterns = rune(pattern[0]), 1
			}
		} else {
			pa, patterns = utf8.RuneError, 0
		}
	}
	return patterns == 0 && stria == 0
}
