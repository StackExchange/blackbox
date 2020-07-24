package bbutil

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAddLinesToSortedFile(t *testing.T) {

	var tests = []struct {
		start    string
		add      []string
		expected string
	}{
		{
			"",
			[]string{"one"},
			"one\n",
		},
		{
			"begin\ntwo\n",
			[]string{"at top"},
			"at top\nbegin\ntwo\n",
		},
		{
			"begin\ntwo\n",
			[]string{"zbottom"},
			"begin\ntwo\nzbottom\n",
		},
		{
			"begin\ntwo\n",
			[]string{"middle"},
			"begin\nmiddle\ntwo\n",
		},
	}

	for i, test := range tests {
		content := []byte(test.start)
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			t.Fatal(err)
		}
		tmpfilename := tmpfile.Name()
		defer os.Remove(tmpfilename)

		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}
		AddLinesToSortedFile(tmpfilename, test.add...)
		expected := test.expected

		got, err := ioutil.ReadFile(tmpfilename)
		if err != nil {
			t.Fatal(err)
		}
		if expected != string(got) {
			t.Errorf("test %v: contents wrong:\nexpected: %q\n     got: %q", i, expected, got)
		}
		os.Remove(tmpfilename)
	}

}
