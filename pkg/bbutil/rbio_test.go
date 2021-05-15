package bbutil

import (
	"testing"
)

func TestRunBashInputOutput(t *testing.T) {

	in := "This is a test of the RBIO system.\n"
	bin := []byte(in)

	out, err := RunBashInputOutput(bin, "cat")
	sout := string(out)
	if err != nil {
		t.Error(err)
	}

	if in != sout {
		t.Errorf("not equal %q %q", in, out)
	}
}
