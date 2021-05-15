package makesafe

// untaint -- A string with a Stringer that is shell safe.

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
	"unicode"
)

type protection int

const (
	// Unknown indicates we don't know if it is safe.
	Unknown protection = iota
	// None requires no special escaping.
	None // Nothing special
	// SingleQuote is unsafe in bash and requires a single quote.
	SingleQuote // Requires at least a single quote
	// DoubleQuote is unsafe in bash and requires escaping or other double-quote features.
	DoubleQuote // Can only be in a double-quoted string
)

const (
	// IsAQuote is either a `'` or `"`
	IsAQuote = None
	// IsSpace is ascii 32
	IsSpace = SingleQuote
	// ShellUnsafe is ()!$ or other bash special char
	ShellUnsafe = SingleQuote
	// GlobUnsafe means could be a glob char (* or ?)
	GlobUnsafe = SingleQuote
	// InterpolationUnsafe used in bash string interpolation ($)
	InterpolationUnsafe = SingleQuote
	// HasBackslash things like \n \t \r \000 \xFF
	HasBackslash = DoubleQuote
)

func max(i, j protection) protection {
	if i > j {
		return i
	}
	return j

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

// Redact returns a string that can be used in a shell single-quoted
// string. It may not be an exact representation, but it is safe
// to include on a command line.
//
// Redacted chars are changed to "X".
// If anything is redacted, the string is surrounded by double quotes
// ("air quotes") and the string "(redacted)" is added to the end.
// If nothing is redacted, but it contains spaces, it is surrounded
// by double quotes.
//
// Example: `s` -> `s`
// Example: `space cadet.txt` -> `"space cadet.txt"`
// Example: `drink a \t soda` -> `"drink a X soda"(redacted)`
// Example: `smile☺` -> `"smile☺`
func Redact(tainted string) string {

	if tainted == "" {
		return `""`
	}

	var b strings.Builder
	b.Grow(len(tainted) + 10)

	redacted := false
	needsQuote := false

	for _, r := range tainted {
		if r == ' ' {
			b.WriteRune(r)
			needsQuote = true
		} else if r == '\'' {
			b.WriteRune('X')
			redacted = true
		} else if r == '"' {
			b.WriteRune('\\')
			b.WriteRune(r)
			needsQuote = true
		} else if unicode.IsPrint(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('X')
			redacted = true
		}
	}

	if redacted {
		return `"` + b.String() + `"(redacted)`
	}
	if needsQuote {
		return `"` + b.String() + `"`
	}
	return tainted
}

// RedactMany returns the list after processing each element with Redact().
func RedactMany(items []string) []string {
	var r []string
	for _, n := range items {
		r = append(r, Redact(n))
	}
	return r
}

// Shell returns the string formatted so that it is safe to be pasted
// into a command line to produce the desired filename as an argument
// to the command.
func Shell(tainted string) string {
	if tainted == "" {
		return `""`
	}

	var b strings.Builder
	b.Grow(len(tainted) + 10)

	level := Unknown
	for _, r := range tainted {
		if r < 128 {
			level = max(level, tab[r].level)
			b.WriteString(tab[r].fn(r))
		} else {
			level = max(level, DoubleQuote)
			b.WriteString(escapeRune(r))
		}
	}
	s := b.String()

	if level == None {
		return tainted
	} else if level == SingleQuote {
		// A single quoted string accepts all chars except the single
		// quote itself, which must be replaced with: '"'"'
		return "'" + strings.Join(strings.Split(s, "'"), `'"'"'`) + "'"
	} else if level == DoubleQuote {
		// A double-quoted string may include \xxx escapes and other
		// things. Sadly bash doesn't interpret those, but printf will!
		return `$(printf '%q' '` + s + `')`
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

// ShellMany returns the list after processing each element with Shell().
func ShellMany(items []string) []string {
	var r []string
	for _, n := range items {
		r = append(r, Redact(n))
	}
	return r
}

// FirstFew returns the first few names. If any are truncated, it is
// noted by appending "...".  The exact definition of "few" may change
// over time, and may be based on the number of chars not the list
func FirstFew(sl []string) string {
	s, _ := FirstFewFlag(sl)
	return s
}

// FirstFewFlag is like FirstFew but returns true if truncation done.
func FirstFewFlag(sl []string) (string, bool) {
	const maxitems = 2
	const maxlen = 70
	if len(sl) < maxitems || len(strings.Join(sl, " ")) < maxlen {
		return strings.Join(sl, " "), false
	}
	return strings.Join(sl[:maxitems], " ") + " (and others)", true
}
