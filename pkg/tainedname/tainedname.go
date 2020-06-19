package tainedname

// tainedname -- A string with a Stringer that is shell safe.

// This goes to great lengths to make sure the String() is pastable.
// Whitespace and shell "special chars" are handled as expected.

// However to be extra paranoid, unicode is turned into backtick
// printf statements.  I don't know anyone that puts unicode in their
// filenames, but I hope they appreciate this.

// Most people would just use strconv.QuoteToGraphic() but I'm a
// control freak.

import (
	"fmt"
	"strings"
)

type dubious string

func New(s string) dubious {
	return dubious(s)
}

type protection int

const (
	Unknown     protection = iota
	None                   // Nothing special
	SingleQuote            // Requires at least a single quote
	DoubleQuote            // Can only be in a double-quoted string
)

var IsAQuote = None                   // Handled as a special case
var IsSpace = SingleQuote             // " " (ascii 32)
var ShellUnsafe = SingleQuote         // bash special
var GlobUnsafe = SingleQuote          // Could be a glob
var InterpolationUnsafe = SingleQuote // Used in bash string interpolation
var HasBackslash = DoubleQuote        // things like \n \t \r \000 \xFF

func max(i, j protection) protection {
	if i > j {
		return i
	} else {
		return j
	}
}

type tabEntry struct {
	level protection
	fn    func(s rune) string
}

var tab [128]tabEntry

func init() {

	for i := 0; i <= 31; i++ { // Control chars
		tab[i] = tabEntry{HasBackslash, oct()}
	}
	tab['\t'] = tabEntry{HasBackslash, literal(`\t`)} // Override
	tab['\n'] = tabEntry{HasBackslash, literal(`\n`)} // Override
	tab['\r'] = tabEntry{HasBackslash, literal(`\r`)} // Override
	tab[' '] = tabEntry{IsSpace, same()}
	tab['!'] = tabEntry{ShellUnsafe, same()}
	tab['"'] = tabEntry{IsAQuote, same()}
	tab['#'] = tabEntry{ShellUnsafe, same()}
	tab['@'] = tabEntry{InterpolationUnsafe, same()}
	tab['$'] = tabEntry{InterpolationUnsafe, same()}
	tab['%'] = tabEntry{InterpolationUnsafe, same()}
	tab['&'] = tabEntry{ShellUnsafe, same()}
	tab['\''] = tabEntry{IsAQuote, same()}
	tab['('] = tabEntry{ShellUnsafe, same()}
	tab[')'] = tabEntry{ShellUnsafe, same()}
	tab['*'] = tabEntry{GlobUnsafe, same()}
	tab['+'] = tabEntry{GlobUnsafe, same()}
	tab[','] = tabEntry{None, same()}
	tab['-'] = tabEntry{None, same()}
	tab['.'] = tabEntry{None, same()}
	tab['/'] = tabEntry{None, same()}
	for i := '0'; i <= '9'; i++ {
		tab[i] = tabEntry{None, same()}
	}
	tab[':'] = tabEntry{InterpolationUnsafe, same()} // ${foo:=default}
	tab[';'] = tabEntry{ShellUnsafe, same()}
	tab['<'] = tabEntry{ShellUnsafe, same()}
	tab['='] = tabEntry{InterpolationUnsafe, same()} // ${foo:=default}
	tab['>'] = tabEntry{ShellUnsafe, same()}
	tab['?'] = tabEntry{GlobUnsafe, same()}
	tab['@'] = tabEntry{InterpolationUnsafe, same()} // ${myarray[@]};
	for i := 'A'; i <= 'Z'; i++ {
		tab[i] = tabEntry{None, same()}
	}
	tab['['] = tabEntry{ShellUnsafe, same()}
	tab['\\'] = tabEntry{ShellUnsafe, same()}
	tab[']'] = tabEntry{GlobUnsafe, same()}
	tab['^'] = tabEntry{GlobUnsafe, same()}
	tab['_'] = tabEntry{None, same()}
	tab['`'] = tabEntry{ShellUnsafe, same()}
	for i := 'a'; i <= 'z'; i++ {
		tab[i] = tabEntry{None, same()}
	}
	tab['{'] = tabEntry{ShellUnsafe, same()}
	tab['|'] = tabEntry{ShellUnsafe, same()}
	tab['}'] = tabEntry{ShellUnsafe, same()}
	tab['~'] = tabEntry{ShellUnsafe, same()}
	tab[127] = tabEntry{HasBackslash, oct()}

	// Check our work. All indexes should have been set.
	for i, e := range tab {
		if e.level == 0 || e.fn == nil {
			panic(fmt.Sprintf("tabEntry %d not set!", i))
		}
	}

}

// literal return this exact string.
func literal(s string) func(s rune) string {
	return func(rune) string { return s }
}

// same converts the rune to a string.
func same() func(r rune) string {
	return func(r rune) string { return string(r) }
}

// oct returns the octal representing the value.
func oct() func(r rune) string {
	return func(r rune) string { return fmt.Sprintf(`\%03o`, r) }
}

// String returns a version of the dirty string that is absolutely
// safe to paste into a command line.
func (dirty dubious) String() string {
	if dirty == "" {
		return `""`
	}

	var b strings.Builder
	b.Grow(len(dirty) + 2)

	level := Unknown
	unicode := false
	for _, r := range dirty {
		if r < 128 {
			level = max(level, tab[r].level)
			b.WriteString(tab[r].fn(r))
		} else {
			level = max(level, DoubleQuote)
			b.WriteString(escapeRune(r))
			unicode = true
		}
	}
	s := b.String()

	switch level {
	case None:
		return string(dirty)
	case SingleQuote:
		// A single quoted string accepts all chars except the single
		// quote itself, which must be replaced with: '"'"'
		return "'" + strings.Join(strings.Split(s, "'"), `'"'"'`) + "'"
	case DoubleQuote:
		if unicode {
			return "$(printf '" + s + "')"
		} else {
			return `"` + s + `"`
		}
	default:
	}
	// should not happen
	return fmt.Sprintf("%q", s)

}

// escapeRune returns a string of octal escapes that represent the rune.
func escapeRune(r rune) string {
	b := []byte(string(rune(r))) // Convert to the indivdual bytes, utf8-encoded.
	// fmt.Printf("rune: len=%d %s %v\n", len(s), s, []byte(s))
	switch len(b) {
	case 1:
		return fmt.Sprintf(`\%03o`, b[0])
	case 2:
		return fmt.Sprintf(`\%03o\%03o`, b[0], b[1])
	case 3:
		return fmt.Sprintf(`\%03o\%03o\%03o`, b[0], b[1], b[2])
	case 4:
		return fmt.Sprintf(`\%03o\%03o\%03o\%03o`, b[0], b[1], b[2], b[3])
	default:
		return string(rune(r))
	}
}
