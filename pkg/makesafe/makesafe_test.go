package makesafe

import (
	"testing"
)

func TestRedact(t *testing.T) {
	for i, test := range []struct{ data, expected string }{
		{"", `""`},
		{"one", "one"},
		{"has space.txt", `"has space.txt"`},
		{"has\ttab.txt", `"hasXtab.txt"(redacted)`},
		{"has\nnl.txt", `"hasXnl.txt"(redacted)`},
		{"has\rret.txt", `"hasXret.txt"(redacted)`},
		{"¡que!", `¡que!`},
		{"thé", `thé`},
		{"pound£", `pound£`},
		{"*.go", `*.go`},
		{"rm -rf / ; echo done", `"rm -rf / ; echo done"`},
		{"smile\u263a", `smile☺`},
		{"dub\U0001D4E6", `dub𝓦`},
		{"four\U0010FFFF", `"fourX"(redacted)`},
	} {
		g := Redact(test.data)
		if g == test.expected {
			t.Logf("%03d: PASSED", i)
		} else {
			t.Errorf("%03d: FAILED data=%q got=(%s) wanted=(%s)", i, test.data, g, test.expected)
		}
	}
}

func TestRedactMany(t *testing.T) {
	data := []string{
		"",
		"one",
		"has space.txt",
		"has\ttab.txt",
	}
	g := RedactMany(data)
	if len(g) != 4 || g[0] != `""` || g[1] != `"has space.txt"` || g[2] != `"hasXtab.txt"(redacted)` {
		t.Logf("PASSED")
	} else {
		t.Errorf("FAILED got=(%q)", g)
	}
}

func TestShell(t *testing.T) {
	for i, test := range []struct{ data, expected string }{
		{"", `""`},
		{"one", "one"},
		{"two\n", `$(printf '%q' 'two\n')`},
		{"ta	tab", `$(printf '%q' 'ta\ttab')`},
		{"tab\ttab", `$(printf '%q' 'tab\ttab')`},
		{"new\nline", `$(printf '%q' 'new\nline')`},
		{"¡que!", `$(printf '%q' '\302\241que!')`},
		{"thé", `$(printf '%q' 'th\303\251')`},
		{"pound£", `$(printf '%q' 'pound\302\243')`},
		{"*.go", `'*.go'`},
		{"rm -rf / ; echo done", `'rm -rf / ; echo done'`},
		{"smile\u263a", `$(printf '%q' 'smile\342\230\272')`},
		{"dub\U0001D4E6", `$(printf '%q' 'dub\360\235\223\246')`},
		{"four\U0010FFFF", `$(printf '%q' 'four\364\217\277\277')`},
	} {
		g := Shell(test.data)
		if g == test.expected {
			t.Logf("%03d: PASSED", i)
			//t.Logf("%03d: PASSED go(%q) bash: %s", i, test.data, test.expected)
		} else {
			t.Errorf("%03d: FAILED data=%q got=`%s` wanted=`%s`", i, test.data, g, test.expected)
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
			t.Logf("%03d: PASSED go=(%q) bash=(%s)", i, test.data, test.expected)
		} else {
			t.Errorf("%03d: FAILED data=%q got=(%s) wanted=(%s)", i, test.data, g, test.expected)
		}
	}
}

func TestShellMany(t *testing.T) {
	data := []string{
		"",
		"one",
		"has space.txt",
		"¡que!",
	}
	g := ShellMany(data)
	if len(g) != 4 || g[0] != `""` || g[1] != "one" || g[2] != `"has space.txt"` || g[3] != `$(printf '%q' '\302\241que!')` {
		t.Logf("PASSED")
	} else {
		t.Errorf("FAILED got=(%q)", g)
	}
}

func TestFirstFewFlag(t *testing.T) {
	for i, test := range []struct {
		data           []string
		expectedFlag   bool
		expectedString string
	}{
		{[]string{"", "one"}, false, ` one`},
		{[]string{"one"}, false, `one`},
		{[]string{"one", "two", "three", "longlonglong", "longlonglonglong", "manylonglonglog", "morelongonglonglong"}, true, ``},
	} {
		gs, gf := FirstFewFlag(test.data)
		if test.expectedFlag {
			if gf == test.expectedFlag {
				t.Logf("%03d: PASSED", i)
			} else {
				t.Errorf("%03d: FAILED data=%q got=(%q) wanted=(%q)", i, test.data, gs, test.expectedString)
			}
		} else {
			if gf == test.expectedFlag && gs == test.expectedString {
				t.Logf("%03d: PASSED", i)
			} else {
				t.Errorf("%03d: FAILED data=%q got=(%q) wanted=(%q)", i, test.data, gs, test.expectedString)
			}
		}
	}
}
