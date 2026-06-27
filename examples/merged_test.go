package capture_test

import (
	"bytes"
	"fmt"

	capture "github.com/gloo-foo/cmd-capture"
	"github.com/gloo-foo/testable"
)

func ExampleCapture_merged() {
	// All stream lines land in a single buffer.
	var merged bytes.Buffer
	_, _ = testable.TestLines(capture.Capture(&merged), "one\ntwo\nthree\n")
	fmt.Print(merged.String())
	// Output:
	// one
	// two
	// three
}
