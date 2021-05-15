package box

import "testing"

func TestPrettyCommitMessage(t *testing.T) {
	long := "aVeryVeryLongLongLongStringStringString"
	for i, test := range []struct {
		data     []string
		expected string
	}{
		{[]string{}, `HEADING (no files)`},
		{[]string{"one"}, `HEADING: one`},
		{[]string{"one", "two"}, `HEADING: one two`},
		{[]string{"one", "two", "three"}, `HEADING: one two three`},
		{[]string{"one", "two", "three", "four"},
			`HEADING: one two three four`},
		{[]string{"one", "two", "three", "four", "five"},
			`HEADING: one two three four five`},
		{[]string{"has spaces.txt"}, `HEADING: "has spaces.txt"`},
		{[]string{"two\n"}, `HEADING: "twoX"(redacted)`},
		{[]string{"smile😁eyes"}, `HEADING: smile😁eyes`},
		{[]string{"tab\ttab", "two very long strings.txt"},
			`HEADING: "tabXtab"(redacted) "two very long strings.txt"`},
		{[]string{long, long, long, long},
			"HEADING: " + long + " " + long + " (and others)"},
	} {
		g := PrettyCommitMessage("HEADING", test.data)
		if g == test.expected {
			//t.Logf("%03d: PASSED files=%q\n", i, test.data)
			t.Logf("%03d: PASSED", i)
		} else {
			t.Errorf("%03d: FAILED files==%q got=(%q) wanted=(%q)\n", i, test.data, g, test.expected)
		}
	}
}
