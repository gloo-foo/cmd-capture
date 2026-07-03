package capture_test

import (
	"bytes"
	"fmt"

	"github.com/gloo-foo/testable"

	capture "github.com/gloo-foo/cmd-capture"
)

func ExampleCapture_basic() {
	// Capture the stream into a buffer while it passes through.
	var buf bytes.Buffer
	_, _ = testable.TestLines(capture.Capture(&buf), "Hello, World!\n")
	fmt.Print(buf.String())
	// Output:
	// Hello, World!
}
