package format

type Style struct {
	Key, String, Number [2]string
	True, False, Null   [2]string
	Append              func(dst []byte, c byte) []byte
}

func hexp(p byte) byte {
	switch {
	case p < 10:
		return p + '0'
	default:
		return (p - 10) + 'a'
	}
}

var TerminalStyle *Style

func init() {
	TerminalStyle = &Style{
		Key:    [2]string{"\x1B[94m", "\x1B[0m"},
		String: [2]string{"\x1B[92m", "\x1B[0m"},
		Number: [2]string{"\x1B[93m", "\x1B[0m"},
		True:   [2]string{"\x1B[96m", "\x1B[0m"},
		False:  [2]string{"\x1B[96m", "\x1B[0m"},
		Null:   [2]string{"\x1B[91m", "\x1B[0m"},
		Append: func(dst []byte, c byte) []byte {
			if c < ' ' && (c != '\r' && c != '\n' && c != '\t' && c != '\v') {
				dst = append(dst, "\\u00"...)
				dst = append(dst, hexp((c>>4)&0xF))
				return append(dst, hexp((c)&0xF))
			}
			return append(dst, c)
		},
	}
}

func Color(src []byte, style *Style) []byte {
	if style == nil {
		style = TerminalStyle
	}
	apnd := style.Append
	if apnd == nil {
		apnd = func(dst []byte, c byte) []byte {
			return append(dst, c)
		}
	}
	type stackt struct {
		kind byte
		key  bool
	}
	var dst []byte
	var stack []stackt
	for i := 0; i < len(src); i++ {
		if src[i] == '"' {
			key := len(stack) > 0 && stack[len(stack)-1].key
			if key {
				dst = append(dst, style.Key[0]...)
			} else {
				dst = append(dst, style.String[0]...)
			}
			dst = apnd(dst, '"')
			for i = i + 1; i < len(src); i++ {
				dst = apnd(dst, src[i])
				if src[i] == '"' {
					j := i - 1
					for ; ; j-- {
						if src[j] != '\\' {
							break
						}
					}
					if (j-i)%2 != 0 {
						break
					}
				}
			}
			if key {
				dst = append(dst, style.Key[1]...)
			} else {
				dst = append(dst, style.String[1]...)
			}
		} else if src[i] == '{' || src[i] == '[' {
			stack = append(stack, stackt{src[i], src[i] == '{'})
			dst = apnd(dst, src[i])
		} else if (src[i] == '}' || src[i] == ']') && len(stack) > 0 {
			stack = stack[:len(stack)-1]
			dst = apnd(dst, src[i])
		} else if (src[i] == ':' || src[i] == ',') && len(stack) > 0 && stack[len(stack)-1].kind == '{' {
			stack[len(stack)-1].key = !stack[len(stack)-1].key
			dst = apnd(dst, src[i])
		} else {
			var kind byte
			if (src[i] >= '0' && src[i] <= '9') || src[i] == '-' {
				kind = '0'
				dst = append(dst, style.Number[0]...)
			} else if src[i] == 't' {
				kind = 't'
				dst = append(dst, style.True[0]...)
			} else if src[i] == 'f' {
				kind = 'f'
				dst = append(dst, style.False[0]...)
			} else if src[i] == 'n' {
				kind = 'n'
				dst = append(dst, style.Null[0]...)
			} else {
				dst = apnd(dst, src[i])
			}
			if kind != 0 {
				for ; i < len(src); i++ {
					if src[i] <= ' ' || src[i] == ',' || src[i] == ':' || src[i] == ']' || src[i] == '}' {
						i--
						break
					}
					dst = apnd(dst, src[i])
				}
				if kind == '0' {
					dst = append(dst, style.Number[1]...)
				} else if kind == 't' {
					dst = append(dst, style.True[1]...)
				} else if kind == 'f' {
					dst = append(dst, style.False[1]...)
				} else if kind == 'n' {
					dst = append(dst, style.Null[1]...)
				}
			}
		}
	}
	return dst
}

