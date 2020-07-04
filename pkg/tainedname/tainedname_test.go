package tainedname

import (
	"testing"
)

func TestRedactUnsafe(t *testing.T) {
	for i, test := range []struct{ data, expected string }{
		{"", `""`},
		{"one", "one"},
		{"has space.txt", "'has space.txt'"},
		{"has\ttab.txt", `hasXtab.txtR`},
		{"has\nnl.txt", `hasXnl.txtR`},
		{"has\rret.txt", `hasXret.txtR`},
		{"¡que!", `¡que!`},
		{"thé", `thé`},
		{"pound£", `pound£`},
		{"*.go", `*.go`},
		{"rm -rf / ; echo done", `'rm -rf / ; echo done'`},
		{"smile\u263a", `smile☺`},
		{"dub\U0001D4E6", `dub𝓦`},
		{"four\U0010FFFF", `fourXR`},
	} {
		g, b := New(test.data).redactHelper()
		if b {
			g = g + "R"
		}
		if g == test.expected {
			//jt.Logf("%03d: PASSED go(%q) bash: %s\n", i, test.data, test.expected)
			t.Logf("%03d: PASSED\n", i)
		} else {
			t.Errorf("%03d: FAILED data=%q got=(%s) wanted=(%s)\n", i, test.data, g, test.expected)
		}
	}
}

func TestString(t *testing.T) {
	for i, test := range []struct{ data, expected string }{
		{"", `""`},
		{"one", "one"},
		{"two\n", `"two\n"`},
		{"tab	tab", `"tab\ttab"`},
		{"tab\ttab", `"tab\ttab"`},
		{"new\nline", `"new\nline"`},
		{"¡que!", `$(printf '\302\241que!')`},
		{"thé", `$(printf 'th\303\251')`},
		{"pound£", `$(printf 'pound\302\243')`},
		{"*.go", `'*.go'`},
		{"rm -rf / ; echo done", `'rm -rf / ; echo done'`},
		{"smile\u263a", `$(printf 'smile\342\230\272')`},
		{"dub\U0001D4E6", `$(printf 'dub\360\235\223\246')`},
		{"four\U0010FFFF", `$(printf 'four\364\217\277\277')`},
	} {
		g := New(test.data).String()
		if g == test.expected {
			t.Logf("%03d: PASSED go(%q) bash: %s\n", i, test.data, test.expected)
		} else {
			t.Errorf("%03d: FAILED data=%q got=(%s) wanted=(%s)\n", i, test.data, g, test.expected)
		}
	}
}

func TestEscapeRune(t *testing.T) {
	for i, test := range []struct {
		data     rune
		expected string
	}{
		{'a', `\141`},
		{'é', `\303\251`},
		{'☺', `\342\230\272`},
		{'글', `\352\270\200`},
		{'𩸽', `\360\251\270\275`},
		//{"\U0010FEDC", `"'\U0010fedc'"`},
	} {
		g := escapeRune(test.data)
		if g == test.expected {
			t.Logf("%03d: PASSED go=(%q) bash=(%s)\n", i, test.data, test.expected)
		} else {
			t.Errorf("%03d: FAILED data=%q got=(%s) wanted=(%s)\n", i, test.data, g, test.expected)
		}
	}
}
