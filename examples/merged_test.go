package capture_test

import (
	"bytes"
	"fmt"

	"github.com/gloo-foo/testable"

	capture "github.com/gloo-foo/cmd-capture"
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
